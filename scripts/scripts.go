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
	filePath := projectPath + genPackageName + "/" + genFileName
	shellPath := files.GetAbPath() + "scripts/comment.sh"
	shellPath = strings.ReplaceAll(shellPath, "!", "\\!")
	op := fmt.Sprintf("bash %s %s %s %s", shellPath, modelDirPath, genPackageName, filePath)
	info, err := cmdc.BashString(op)
	fmt.Println(info)
	return err
}

func RunTestReport(moduleName string) error {
	shellPath := files.GetAbPath() + "scripts/test.sh"
	shellPath = strings.ReplaceAll(shellPath, "!", "\\!")
	op := fmt.Sprintf("bash %s %s", shellPath, moduleName)
	info, err := cmdc.BashString(op)
	fmt.Println(info)
	return err
}
