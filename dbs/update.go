package dbs

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/YouAreOnlyOne/goutils/files"
	"github.com/YouAreOnlyOne/goutils/times"
	"io"
	"os"
	"strings"
	"time"
)

type UpdateModel struct {
	BaseModel
	FileName    string    `json:"file_name" gorm:"column:file_name;not null;comment:'文件名称'"`
	ExecuteTime time.Time `json:"execute_time" gorm:"column:execute_time;null;comment:'执行时间'"`
}

func (m *UpdateModel) TableName() string {
	return "update_info"
}

func (db *DB) InitDefaultData(basePath string) {
	db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&UpdateModel{})
	fileList := files.GetPathBySuffix(basePath, ".sql")
	recordList := make([]*UpdateModel, 0)
	err := db.Find(&recordList).Error
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var recordMap = make(map[string]int)
	for _, record := range recordList {
		recordMap[record.FileName] = 1
	}

	count := 0
	var buffer bytes.Buffer
	sql := "insert into `update_info` (`file_name`,`execute_time`) values "
	buffer.WriteString(sql)
	now := time.Now().Format(times.TmFmtWithMS1)
	for i, file := range fileList {
		fileName := file[strings.LastIndex(file, "/")+1:]
		if recordMap[fileName] > 0 {
			continue
		}
		count++
		if i < len(fileList)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s'),", fileName, now))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s');", fileName, now))
		}
		f, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()
		reader := bufio.NewReader(f)
		var lines []string
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			content := strings.TrimSpace(string(line))
			if content == "" || strings.HasPrefix(content, "--") {
				continue
			}
			lines = append(lines, content)
		}
		sqlContent := strings.Join(lines, "")
		sqlList := strings.Split(sqlContent, ";")
		for _, sql := range sqlList {
			if sql == "" {
				continue
			}
			err = db.DB.Exec(sql).Error
			if err != nil {
				fmt.Println(sql, err)
				os.Exit(1)
			}
		}

	}
	if count > 0 {
		err = db.DB.Exec(buffer.String()).Error
		if err != nil {
			fmt.Println(buffer.String(), err)
			os.Exit(1)
		}
	}

}
