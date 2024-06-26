/**
* @author: yangchangjia
* @email 1320259466@qq.com
* @date: 2024/4/17 14:55
* @desc: about the role of class.
 */

package domain

// Package generator downloads an updated version of the PSL list and compiles it into go code.
//
// It is meant to be used by maintainers in conjunction with the go generate tool
// to update the list.

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-github/v50/github"
	"go/format"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

const (
	list = `// This file is automatically generated
// Run "go run domain/gen.go" to update the list.

package domain

const ListVersion = "PSL version {{.VersionSHA}} ({{.VersionDate}})"

func DefaultRules() [{{len .Rules}}]Rule {
	return r
}

var r = [{{len .Rules}}]Rule{
	{{range $r := .Rules}} \
	{ {{$r.Type}}, "{{$r.Value}}", {{$r.Length}}, {{$r.Private}} },
	{{end}}
}

func init() {
	for i := range r {
		DefaultList.AddRule(&r[i])
	}
}

`
)

var (
	listTmpl = template.Must(template.New("list").Parse(cont(list)))
)

// https://github.com/golang/go/issues/9969
// Requires go1.6
func cont(s string) string {
	return strings.Replace(s, "\\\n", "", -1)
}

func extractHeadInfo() (sha string, datetime github.Timestamp) {
	client := github.NewClient(nil)

	commits, _, err := client.Repositories.ListCommits(context.Background(), "publicsuffix", "list", nil)
	if err != nil {
		log.Fatal(err)
	}

	lastCommit := commits[0]
	return lastCommit.GetSHA(), lastCommit.GetCommit().GetCommitter().GetDate()
}

// Generator represents a generator.
type Generator struct {
	Verbose bool
}

// NewGenerator creates a Generator with default settings.
func NewGenerator() *Generator {
	g := &Generator{
		Verbose: false,
	}
	return g
}

// Write ...
func (g *Generator) Write(filename string) error {
	content, err := g.generate()
	if err != nil {
		return err
	}

	g.log("Writing %v...\n", filename)
	return ioutil.WriteFile(filename, content, 0644)
}

// Print ...
func (g *Generator) Print() error {
	content, err := g.generate()
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(content)
	return err
}

// Generate downloads an updated version of the PSL list and compiles it into go code.
func (g *Generator) generate() ([]byte, error) {
	g.log("Fetching PSL version...\n")
	sha, datetime := extractHeadInfo()

	g.log("Downloading PSL %s...\n", sha[:6])
	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/publicsuffix/list/%s/public_suffix_list.dat", sha))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	list := NewList()
	rules, err := list.Load(resp.Body, nil)
	if err != nil {
		return nil, err
	}

	data := struct {
		VersionSHA  string
		VersionDate string
		Rules       []Rule
	}{
		sha[:6],
		datetime.Format(time.ANSIC),
		rules,
	}

	g.log("Parsing PSL...\n")
	buf := new(bytes.Buffer)
	err = listTmpl.Execute(buf, &data)
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

func (g *Generator) log(format string, v ...interface{}) {
	if !g.Verbose {
		return
	}

	log.Printf(format, v...)
}

func (g *Generator) fatal(message string) {
	if !g.Verbose {
		fmt.Println(message)
		os.Exit(1)
	}

	log.Fatal(message)
}
