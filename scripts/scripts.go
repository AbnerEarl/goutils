/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/4/16 11:46
 * @desc: about the role of class.
 */

package scripts

import (
	"fmt"
	"github.com/AbnerEarl/goutils/cmdc"
	"github.com/AbnerEarl/goutils/files"
	"strings"
)

func GenDbComment(modelDirPath, genPackageName, genFileName string) error {
	projectPath := files.GetParentPath(modelDirPath)
	filePath := projectPath + genFileName
	shellPath := files.GetAbPathByCaller() + "scripts/comment.sh"
	shellPath = strings.ReplaceAll(shellPath, "!", "\\!")
	op := fmt.Sprintf("bash %s %s %s %s", shellPath, modelDirPath, genPackageName, filePath)
	return cmdc.Bash(op)
}

func RunTestReport(moduleName string) error {
	shellPath := files.GetAbPathByCaller() + "scripts/test.sh"
	shellPath = strings.ReplaceAll(shellPath, "!", "\\!")
	op := fmt.Sprintf("bash %s %s", shellPath, moduleName)
	return cmdc.Bash(op)
}
