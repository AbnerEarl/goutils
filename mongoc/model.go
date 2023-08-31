package mongoc

import (
	"bufio"
	"fmt"
	"github.com/AbanerEarl/goutils/files"
	"io"
	"os"
	"strings"
	"time"
)

type BaseModel struct {
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"-" bson:"deleted_at"`
	Remark    string     `json:"remark" bson:"remark"`
}

func (c *BaseModel) CollectionName() string {
	return "base_info"
}

type UpdateModel struct {
	BaseModel
	FileName    string    `json:"file_name" bson:"file_name"`
	ExecuteTime time.Time `json:"execute_time" bson:"execute_time"`
}

func (c *UpdateModel) CollectionName() string {
	return "update_info"
}

func InitDefaultData(basePath string) {
	up := UpdateModel{}
	fileList := files.GetPathBySuffix(basePath, ".sql")
	recordList := make([]*UpdateModel, 0)
	err := FindMany(DefaultMaxLimit, 0, up.CollectionName(), nil, nil, &recordList, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var recordMap = make(map[string]int)
	for _, record := range recordList {
		recordMap[record.FileName] = 1
	}

	count := 0
	var addList []interface{}
	for _, file := range fileList {
		fileName := file[strings.LastIndex(file, "/")+1:]
		if recordMap[fileName] > 0 {
			continue
		}
		count++
		addList = append(addList, UpdateModel{
			FileName:    fileName,
			ExecuteTime: time.Now(),
		})
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
		sqlContent = strings.ReplaceAll(sqlContent, "\n", "")
		sqlList := strings.Split(sqlContent, ";")
		for _, sql := range sqlList {
			if len(sql) < 10 {
				continue
			}
			sql = strings.ReplaceAll(sql, "`", "")
			sql = strings.ReplaceAll(sql, "'", "")
			coll := strings.Split(sql, " ")[2]
			names := strings.Split(sql[strings.Index(sql, "("):strings.Index(sql, ")")], ",")
			sql = strings.Split(sql, "values")[1]
			values := strings.Split(sql[strings.Index(sql, "(")+1:strings.LastIndex(sql, ")")], ",")
			m := map[string]interface{}{}
			for i, n := range names {
				if strings.Contains(values[i], "now(") {
					m[n] = time.Now()
				} else {
					m[n] = values[i]
				}

			}
			err = InsertOne(coll, m)
			if err != nil {
				fmt.Println(sql, err)
				os.Exit(1)
			}
		}

	}
	if count > 0 {
		err = InsertMany(up.CollectionName(), addList)
		if err != nil {
			fmt.Println(addList, err)
			os.Exit(1)
		}
	}

}
