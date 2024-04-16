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

func GenSwagDoc(apiPathList []string) error {
	cmd := "swag init --parseDependency --parseInternal"
	if len(apiPathList) > 0 {
		cmd = fmt.Sprintf("%s --dir %s", cmd, strings.Join(apiPathList, ","))
	}
	return cmdc.Bash(cmd)
}

func InstallSwag() error {
	return cmdc.Bash("go install github.com/swaggo/swag/cmd/swag@v1.16.3")
}
