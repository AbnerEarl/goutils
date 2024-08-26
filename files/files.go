/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/4/3 14:45
 * @desc: about the role of class.
 */

package files

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetAbPath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(path.Dir(filename)) + "/"
	}
	return abPath
}

func GetCuPath() string {
	cu, _ := os.Getwd()
	abp, _ := filepath.Abs(cu)
	return abp
}

func GetParentPath(srcPath string) string {
	for strings.HasSuffix(srcPath, "/") {
		srcPath = srcPath[:len(srcPath)-1]
	}
	return srcPath[:strings.LastIndex(srcPath, "/")+1]
}

func GetFilesBySuffix(path string, suffix string) []string {
	var allFiles []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return allFiles
	}

	for _, file := range files {
		if file.IsDir() {
			allFiles = append(allFiles, GetFilesBySuffix(path+"/"+file.Name(), suffix)...)
		} else {
			if strings.HasSuffix(file.Name(), suffix) {
				allFiles = append(allFiles, path+"/"+file.Name())
			}
		}
	}
	return allFiles
}

func MoveChildToParent(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			childPath := dirPath + string(filepath.Separator) + file.Name()
			childFiles, err := ioutil.ReadDir(childPath)
			if err != nil {
				return err
			}
			for _, f := range childFiles {
				err := os.Rename(childPath+string(filepath.Separator)+f.Name(), dirPath+string(filepath.Separator)+f.Name())
				if err != nil {
					return err
				}
			}
			os.RemoveAll(childPath)
		}
	}
	return nil
}

func CompressFileName(filePath string) string {
	if strings.HasSuffix(filePath, ".tar.gz") {
		return filePath[:strings.LastIndex(filePath, ".tar.gz")]
	} else if strings.HasSuffix(filePath, ".tar") {
		return filePath[:strings.LastIndex(filePath, ".tar")]
	} else if strings.HasSuffix(filePath, ".tar.bz2") {
		return filePath[:strings.LastIndex(filePath, ".tar.bz2")]
	} else if strings.HasSuffix(filePath, ".tar.Z") {
		return filePath[:strings.LastIndex(filePath, ".tar.Z")]
	} else if strings.HasSuffix(filePath, ".zip") {
		return filePath[:strings.LastIndex(filePath, ".zip")]
	}
	return ""
}

func DecompressFile(filePath, destDir string) error {
	var op string
	if strings.HasSuffix(filePath, ".tar.gz") {
		op = fmt.Sprintf("tar -zxvf %s -C %s", filePath, destDir)
	} else if strings.HasSuffix(filePath, ".tar") {
		op = fmt.Sprintf("tar -xvf %s -C %s", filePath, destDir)
	} else if strings.HasSuffix(filePath, ".tar.bz2") {
		op = fmt.Sprintf("tar -xjvf %s -C %s", filePath, destDir)
	} else if strings.HasSuffix(filePath, ".tar.Z") {
		op = fmt.Sprintf("tar -zxvf %s -C %s", filePath, destDir)
	} else if strings.HasSuffix(filePath, ".zip") {
		op = fmt.Sprintf("unzip -o %s -d %s", filePath, destDir)
	} else {
		return fmt.Errorf("not surport file format")
	}
	cmd := exec.Command("sh", "-c", op)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return err
}
