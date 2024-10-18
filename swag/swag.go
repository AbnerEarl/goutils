/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/4/16 14:54
 * @desc: about the role of class.
 */

package swag

import (
	"fmt"
	"github.com/AbnerEarl/goutils/cmdc"
	"strings"
)

func GenSwagDoc(mainFile, projectPath string, apiDirPaths []string) error {
	cmd := "swag init --parseDependency --parseInternal"
	if len(projectPath) > 0 {
		cmd = fmt.Sprintf("%s --dir %s", cmd, projectPath)
	} else {
		projectPath = "."
	}
	if len(mainFile) > 0 {
		cmd = fmt.Sprintf("%s -g %s", cmd, mainFile)
	}
	if len(apiDirPaths) > 0 {
		cmd = fmt.Sprintf("%s --exclude $(find '%s' -type d -maxdepth 1  -not -name '%s' | grep -vE '%s' | tr '\\n' ',')", cmd, projectPath, projectPath, strings.Join(apiDirPaths, "|"))
	}
	fmt.Println(cmd)
	info, err := cmdc.BashString(cmd)
	fmt.Println(info)
	return err
}

func InstallSwag() error {
	//info, err := cmdc.BashString("go install github.com/swaggo/swag/cmd/swag@v1.8.3")
	info, err := cmdc.BashString("go install github.com/swaggo/swag/cmd/swag@v1.8.12")
	fmt.Println(info)
	return err
}
