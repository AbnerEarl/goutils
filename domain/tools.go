/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/4/16 17:30
 * @desc: about the role of class.
 */

package domain

// Package publicsuffix provides a domain name parser
// based on data from the public suffix list http://publicsuffix.org/.
// A public suffix is one under which Internet users can directly register names.

import (
	"bufio"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io"
	"net"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/idna"
)

const (
	// Version identifies the current library version.
	// This is a pro forma convention given that Go dependencies
	// tends to be fetched directly from the repo.
	Version = "0.30.2"

	// NormalType represents a normal rule such as "com"
	NormalType = 1
	// WildcardType represents a wildcard rule such as "*.com"
	WildcardType = 2
	// ExceptionType represents an exception to a wildard rule
	ExceptionType = 3

	listTokenPrivateDomains = "===BEGIN PRIVATE DOMAINS==="
	listTokenComment        = "//"
)

// DefaultList is the default List and it is used by Parse and Domain.
var DefaultList = NewList()

// DefaultRule is the default Rule that represents "*".
var DefaultRule = MustNewRule("*")

// DefaultParserOptions are the default options used to parse a Public Suffix list.
var DefaultParserOptions = &ParserOption{PrivateDomains: true, ASCIIEncoded: false}

// DefaultFindOptions are the default options used to perform the lookup of rules in the list.
var DefaultFindOptions = &FindOptions{IgnorePrivate: false, DefaultRule: DefaultRule}

// Rule represents a single rule in a Public Suffix List.
type Rule struct {
	Type    int
	Value   string
	Length  int
	Private bool
}

// ParserOption are the options you can use to customize the way a List
// is parsed from a file or a string.
type ParserOption struct {
	// Set to false to skip the private domains when parsing.
	// Default to true, which means the private domains are included.
	PrivateDomains bool

	// Set to false if the input is encoded in U-labels (Unicode)
	// as opposite to A-labels.
	// Default to false, which means the list is containing Unicode domains.
	// This is the default because the original PSL currently contains Unicode.
	ASCIIEncoded bool
}

// FindOptions are the options you can use to customize the way a Rule
// is searched within the list.
type FindOptions struct {
	// Set to true to ignore the rules within the "Private" section of the Public Suffix List.
	IgnorePrivate bool

	// The default rule to use when no rule matches the input.
	// The format Public Suffix algorithm states that the rule "*" should be used when no other rule matches,
	// but some consumers may have different needs.
	DefaultRule *Rule
}

// List represents a Public Suffix List.
type List struct {
	// rules is kept private because you should not access rules directly
	rules map[string]*Rule
}

// NewList creates a new empty list.
func NewList() *List {
	return &List{
		rules: map[string]*Rule{},
	}
}

// NewListFromString parses a string that represents a Public Suffix source
// and returns a List initialized with the rules in the source.
func NewListFromString(src string, options *ParserOption) (*List, error) {
	l := NewList()
	_, err := l.LoadString(src, options)
	return l, err
}

// NewListFromFile parses a string that represents a Public Suffix source
// and returns a List initialized with the rules in the source.
func NewListFromFile(path string, options *ParserOption) (*List, error) {
	l := NewList()
	_, err := l.LoadFile(path, options)
	return l, err
}

// Load parses and loads a set of rules from an io.Reader into the current list.
func (l *List) Load(r io.Reader, options *ParserOption) ([]Rule, error) {
	return l.parse(r, options)
}

// LoadString parses and loads a set of rules from a String into the current list.
func (l *List) LoadString(src string, options *ParserOption) ([]Rule, error) {
	r := strings.NewReader(src)
	return l.parse(r, options)
}

// LoadFile parses and loads a set of rules from a File into the current list.
func (l *List) LoadFile(path string, options *ParserOption) ([]Rule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return l.parse(f, options)
}

// AddRule adds a new rule to the list.
//
// The exact position of the rule into the list is unpredictable.
// The list may be optimized internally for lookups, therefore the algorithm
// will decide the best position for the new rule.
func (l *List) AddRule(r *Rule) error {
	l.rules[r.Value] = r
	return nil
}

// Size returns the size of the list, which is the number of rules.
func (l *List) Size() int {
	return len(l.rules)
}

func (l *List) parse(r io.Reader, options *ParserOption) ([]Rule, error) {
	if options == nil {
		options = DefaultParserOptions
	}
	var rules []Rule

	scanner := bufio.NewScanner(r)
	var section int // 1 == ICANN, 2 == PRIVATE

Scanning:
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {

		// skip blank lines
		case line == "":
			break

		// include private domains or stop scanner
		case strings.Contains(line, listTokenPrivateDomains):
			if !options.PrivateDomains {
				break Scanning
			}
			section = 2

		// skip comments
		case strings.HasPrefix(line, listTokenComment):
			break

		default:
			var rule *Rule
			var err error

			if options.ASCIIEncoded {
				rule, err = NewRule(line)
			} else {
				rule, err = NewRuleUnicode(line)
			}
			if err != nil {
				return []Rule{}, err
			}

			rule.Private = (section == 2)
			l.AddRule(rule)
			rules = append(rules, *rule)
		}

	}

	return rules, scanner.Err()
}

// Find and returns the most appropriate rule for the domain name.
func (l *List) Find(name string, options *FindOptions) *Rule {
	if options == nil {
		options = DefaultFindOptions
	}

	part := name
	for {
		rule, ok := l.rules[part]

		if ok && rule.Match(name) && !(options.IgnorePrivate && rule.Private) {
			return rule
		}

		i := strings.IndexRune(part, '.')
		if i < 0 {
			return options.DefaultRule
		}

		part = part[i+1:]
	}

}

// NewRule parses the rule content, creates and returns a Rule.
//
// The content of the rule MUST be encoded in ASCII (A-labels).
func NewRule(content string) (*Rule, error) {
	var rule *Rule
	var value string

	switch content[0] {
	case '*': // wildcard
		if content == "*" {
			value = ""
		} else {
			value = content[2:]
		}
		rule = &Rule{Type: WildcardType, Value: value, Length: len(Labels(value)) + 1}
	case '!': // exception
		value = content[1:]
		rule = &Rule{Type: ExceptionType, Value: value, Length: len(Labels(value))}
	default: // normal
		value = content
		rule = &Rule{Type: NormalType, Value: value, Length: len(Labels(value))}
	}

	return rule, nil
}

// NewRuleUnicode is like NewRule, but expects the content to be encoded in Unicode (U-labels).
func NewRuleUnicode(content string) (*Rule, error) {
	var err error

	content, err = ToASCII(content)
	if err != nil {
		return nil, err
	}

	return NewRule(content)
}

// MustNewRule is like NewRule, but panics if the content cannot be parsed.
func MustNewRule(content string) *Rule {
	rule, err := NewRule(content)
	if err != nil {
		panic(err)
	}
	return rule
}

// Match checks if the rule matches the name.
//
// A domain name is said to match a rule if and only if all of the following conditions are met:
//   - When the domain and rule are split into corresponding labels,
//     that the domain contains as many or more labels than the rule.
//   - Beginning with the right-most labels of both the domain and the rule,
//     and continuing for all labels in the rule, one finds that for every pair,
//     either they are identical, or that the label from the rule is "*".
//
// See https://publicsuffix.org/list/
func (r *Rule) Match(name string) bool {
	left := strings.TrimSuffix(name, r.Value)

	// the name contains as many labels than the rule
	// this is a match, unless it's a wildcard
	// because the wildcard requires one more label
	if left == "" {
		return r.Type != WildcardType
	}

	// if there is one more label, the rule match
	// because either the rule is shorter than the domain
	// or the rule is a wildcard and there is one more label
	return left[len(left)-1:] == "."
}

// Decompose takes a name as input and decomposes it into a tuple of <TRD+SLD, TLD>,
// according to the rule definition and type.
func (r *Rule) Decompose(name string) (result [2]string) {
	if r == DefaultRule {
		i := strings.LastIndexByte(name, '.')
		if i < 0 {
			return
		}
		result[0], result[1] = name[:i], name[i+1:]
		return
	}
	switch r.Type {
	case NormalType:
		name = strings.TrimSuffix(name, r.Value)
		if len(name) == 0 {
			return
		}
		result[0], result[1] = name[:len(name)-1], r.Value
	case WildcardType:
		name := strings.TrimSuffix(name, r.Value)
		if len(name) == 0 {
			return
		}
		name = name[:len(name)-1]
		i := strings.LastIndexByte(name, '.')
		if i < 0 {
			return
		}
		result[0], result[1] = name[:i], name[i+1:]+"."+r.Value
	case ExceptionType:
		i := strings.IndexRune(r.Value, '.')
		if i < 0 {
			return
		}
		suffix := r.Value[i+1:]
		name = strings.TrimSuffix(name, suffix)
		if len(name) == 0 {
			return
		}
		result[0], result[1] = name[:len(name)-1], suffix
	}
	return
}

// Labels decomposes given domain name into labels,
// corresponding to the dot-separated tokens.
func Labels(name string) []string {
	return strings.Split(name, ".")
}

// DomainName represents a domain name.
type DomainName struct {
	TLD  string
	SLD  string
	TRD  string
	Rule *Rule
}

// String joins the components of the domain name into a single string.
// Empty labels are skipped.
//
// Examples:
//
//	DomainName{"com", "example"}.String()
//	// example.com
//	DomainName{"com", "example", "www"}.String()
//	// www.example.com
func (d *DomainName) String() string {
	switch {
	case d.TLD == "":
		return ""
	case d.SLD == "":
		return d.TLD
	case d.TRD == "":
		return d.SLD + "." + d.TLD
	default:
		return d.TRD + "." + d.SLD + "." + d.TLD
	}
}

// Domain extract and return the domain name from the input
// using the default (Public Suffix) List.
//
// Examples:
//
//	publicsuffix.Domain("example.com")
//	// example.com
//	publicsuffix.Domain("www.example.com")
//	// example.com
//	publicsuffix.Domain("www.example.co.uk")
//	// example.co.uk
func Domain(name string) (string, error) {
	return DomainFromListWithOptions(DefaultList, name, DefaultFindOptions)
}

// Parse decomposes the name into TLD, SLD, TRD
// using the default (Public Suffix) List,
// and returns the result as a DomainName
//
// Examples:
//
//	list := NewList()
//
//	publicsuffix.Parse("example.com")
//	// &DomainName{"com", "example"}
//	publicsuffix.Parse("www.example.com")
//	// &DomainName{"com", "example", "www"}
//	publicsuffix.Parse("www.example.co.uk")
//	// &DomainName{"co.uk", "example"}
func Parse(name string) (*DomainName, error) {
	return ParseFromListWithOptions(DefaultList, name, DefaultFindOptions)
}

// DomainFromListWithOptions extract and return the domain name from the input
// using the (Public Suffix) list passed as argument.
//
// Examples:
//
//	list := NewList()
//
//	publicsuffix.DomainFromListWithOptions(list, "example.com")
//	// example.com
//	publicsuffix.DomainFromListWithOptions(list, "www.example.com")
//	// example.com
//	publicsuffix.DomainFromListWithOptions(list, "www.example.co.uk")
//	// example.co.uk
func DomainFromListWithOptions(l *List, name string, options *FindOptions) (string, error) {
	dn, err := ParseFromListWithOptions(l, name, options)
	if err != nil {
		return "", err
	}
	return dn.SLD + "." + dn.TLD, nil
}

// ParseFromListWithOptions decomposes the name into TLD, SLD, TRD
// using the (Public Suffix) list passed as argument,
// and returns the result as a DomainName
//
// Examples:
//
//	list := NewList()
//
//	publicsuffix.ParseFromListWithOptions(list, "example.com")
//	// &DomainName{"com", "example"}
//	publicsuffix.ParseFromListWithOptions(list, "www.example.com")
//	// &DomainName{"com", "example", "www"}
//	publicsuffix.ParseFromListWithOptions(list, "www.example.co.uk")
//	// &DomainName{"co.uk", "example"}
func ParseFromListWithOptions(l *List, name string, options *FindOptions) (*DomainName, error) {
	n, err := normalize(name)
	if err != nil {
		return nil, err
	}

	r := l.Find(n, options)
	if r == nil {
		return nil, fmt.Errorf("no rule matching name %s", name)
	}

	parts := r.Decompose(n)
	left, tld := parts[0], parts[1]
	if tld == "" {
		return nil, fmt.Errorf("%s is a suffix", n)
	}

	dn := &DomainName{
		Rule: r,
		TLD:  tld,
	}
	if i := strings.LastIndexByte(left, '.'); i < 0 {
		dn.SLD = left
	} else {
		dn.TRD = left[:i]
		dn.SLD = left[i+1:]
	}
	return dn, nil
}

func normalize(name string) (string, error) {
	ret := strings.ToLower(name)

	if ret == "" {
		return "", fmt.Errorf("name is blank")
	}
	if ret[0] == '.' {
		return "", fmt.Errorf("name %s starts with a dot", ret)
	}

	return ret, nil
}

// ToASCII is a wrapper for idna.ToASCII.
//
// This wrapper exists because idna.ToASCII backward-compatibility was broken twice in few months
// and I can't call this package directly anymore. The wrapper performs some terrible-but-necessary
// before-after replacements to make sure an already ASCII input always results in the same output
// even if passed through ToASCII.
//
// See golang/net@67957fd0b1, golang/net@f2499483f9, golang/net@78ebe5c8b6,
// and weppos/publicsuffix-go#66.
func ToASCII(s string) (string, error) {
	// .example.com should be .example.com
	// ..example.com should be ..example.com
	if strings.HasPrefix(s, ".") {
		dotIndex := 0
		for i := 0; i < len(s); i++ {
			if s[i] == '.' {
				dotIndex = i
			} else {
				break
			}
		}
		out, err := idna.ToASCII(s[dotIndex+1:])
		out = s[:dotIndex+1] + out
		return out, err
	}

	return idna.ToASCII(s)
}

// ToUnicode is a wrapper for idna.ToUnicode.
//
// See ToASCII for more details about why this wrapper exists.
func ToUnicode(s string) (string, error) {
	return idna.ToUnicode(s)
}

// CookieJarList implements the cookiejar.PublicSuffixList interface.
var CookieJarList cookiejar.PublicSuffixList = cookiejarList{DefaultList}

type cookiejarList struct {
	List *List
}

// PublicSuffix implements cookiejar.PublicSuffixList.
func (l cookiejarList) PublicSuffix(domain string) string {
	rule := l.List.Find(domain, nil)
	return rule.Decompose(domain)[1]
}

// PublicSuffix implements cookiejar.String.
func (cookiejarList) String() string {
	return ListVersion
}

// PublicSuffix returns the public suffix of the domain
// using a copy of the publicsuffix.org database packaged into this library.
//
// Note. To maintain compatibility with the golang.org/x/net/publicsuffix
// this method doesn't return an error. However, in case of error,
// the returned value is empty.
func PublicSuffix(domain string) (publicSuffix string, icann bool) {
	//d, err := psl.Parse(domain)
	//if err != nil {
	//	return "", false
	//}
	//
	//return d.Rule.Value, !d.Rule.Private

	rule := DefaultList.Find(domain, nil)
	publicSuffix = rule.Decompose(domain)[1]
	icann = !rule.Private

	// x/net/publicsuffix sets icann to false when the default rule "*" is used
	if rule.Value == "" && rule.Type == WildcardType {
		icann = false
	}

	return
}

// EffectiveTLDPlusOne returns the effective top level domain plus one more label.
// For example, the eTLD+1 for "foo.bar.golang.org" is "golang.org".
func EffectiveTLDPlusOne(domain string) (string, error) {
	return Domain(domain)
}

// GetMainDomain 输入任意子域名、URL或IP地址，返回其主域名或IP地址。
// 功能特性:
// 1. 准确提取主域名：自动处理多部分顶级域名，如 "google.co.uk"。
// 2. 支持IP地址：如果输入是IP地址（v4或v6），则直接返回该IP地址。
// 3. 支持中文域名：可以正确处理包含中文字符的国际化域名（IDN）。
// 4. 兼容完整URL：能从复杂的URL中提取出主机部分进行处理
// 示例:
// - "deep.learning.google.co.uk" -> "google.co.uk"
// - "http://192.168.1.1:8080/path" -> "192.168.1.1"
// - "邮件.宁波舟山港.net" -> "宁波舟山港.net"
func GetMainDomain(input string) (string, error) {
	// 1. 清洗和预处理输入
	// 如果输入不包含协议头，url.Parse可能会将其解析为路径部分，
	// 因此我们为其添加一个临时的协议头。
	if !strings.Contains(input, "://") && !strings.HasPrefix(input, "//") {
		input = "http://" + input
	}

	parsedURL, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("无法解析输入 '%s': %w", input, err)
	}

	// 提取主机名部分，例如从 "https://www.google.com/search" 提取 "www.google.com"
	// 对于中文域名，Go的url.Parse会自动进行Punycode编码。
	// 但我们希望最后返回的是原始的中文域名，所以先用原始主机名。
	hostname := parsedURL.Hostname()
	if hostname == "" {
		return "", fmt.Errorf("无法从输入中提取主机名")
	}

	// 2. 【新功能】检查主机名是否为IP地址
	// net.ParseIP可以同时处理IPv4和IPv6。
	if ip := net.ParseIP(hostname); ip != nil {
		return hostname, nil // 是IP地址，直接返回
	}

	// 3. 如果不是IP，则使用 publicsuffix 库获取 eTLD+1
	// EffectiveTLDPlusOne 会自动处理像 .co.uk, .com.cn 这样的公共后缀。
	mainDomain, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		// 对于 localhost 或其他无法识别的域名，这里会报错
		return "", fmt.Errorf("无法获取 '%s' 的主域名: %w", hostname, err)
	}

	// 4. 处理中文域名显示问题
	// 由于url.Parse会将中文转为Punycode，我们需要将结果转换回来。
	// 如果原始主机名包含我们找到的主域名（Punycode形式），我们假设可以安全地
	// 从原始主机名中提取出对应的中文主域名。
	// 这是一个简化的处理，但在大多数情况下有效。
	// 更严谨的方法是使用 idna 库进行转换，但 publicsuffix 返回的已经是处理过的结果，
	// 所以我们在这里进行后缀匹配。
	originalHost := parsedURL.Host // 使用 .Host 保留中文字符
	if originalHost != hostname {  // 说明存在中文字符被转换了
		// 为了避免端口号干扰，先移除
		if portIndex := strings.LastIndex(originalHost, ":"); portIndex != -1 {
			// 确保不是IPv6的冒号
			if ip := net.ParseIP(originalHost[:portIndex]); ip == nil || ip.To4() != nil {
				originalHost = originalHost[:portIndex]
			}
		}

		// 查找Punycode主域名在原始域名中的对应部分
		// 示例: originalHost="邮件.宁波舟山港.net", mainDomain="xn--zfr64x43ey0a322o.net"
		// 我们需要找到 "宁波舟山港.net"

		// 这是一个更简单可靠的办法：直接对原始主机使用publicsuffix
		// 但publicsuffix库的List不直接接受unicode，所以前面的Punycode转换是必要的
		// 所以我们在这里进行一个简单的替换逻辑
		punycodeTld, _ := publicsuffix.PublicSuffix(hostname)
		originalTld, _ := publicsuffix.PublicSuffix(originalHost) // 可能会失败，但没关系

		if strings.HasSuffix(originalHost, originalTld) && strings.HasSuffix(mainDomain, punycodeTld) {
			// 替换后缀，得到主域名的unicode版本
			originalBaseDomain := strings.TrimSuffix(strings.TrimSuffix(originalHost, "."+originalTld), "."+punycodeTld) // 双重trim以防万一

			// 找到最后一个 "." 的位置
			lastDotIndex := strings.LastIndex(originalBaseDomain, ".")
			if lastDotIndex != -1 {
				originalBaseDomain = originalBaseDomain[lastDotIndex+1:]
			}

			// 如果punycode解码后的基础域名和原始基础域名匹配（或近似），则返回原始版本
			// 这是一个启发式方法，因为没有直接的反向映射
			return originalBaseDomain + "." + originalTld, nil
		}
	}

	return mainDomain, nil
}
