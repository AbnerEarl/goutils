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

func GenSwagDoc(mainDirPath string, apiDirPaths []string) error {
	cmd := "swag init --parseDependency --parseInternal"
	if len(mainDirPath) > 0 {
		cmd = fmt.Sprintf("%s --dir %s", cmd, mainDirPath)
	} else {
		mainDirPath = "."
	}
	if len(apiDirPaths) > 0 {
		cmd = fmt.Sprintf("%s --exclude $(find '%s' -type d -maxdepth 1  -not -name '%s' | grep -vE '%s' | tr '\\n' ',')", cmd, mainDirPath, mainDirPath, strings.Join(apiDirPaths, "|"))
	}
	return cmdc.Bash(cmd)
}

func InstallSwag() error {
	return cmdc.Bash("go install github.com/swaggo/swag/cmd/swag@v1.8.3")
}
