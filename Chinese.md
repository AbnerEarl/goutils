# goutils
一些非常有用的中间件库二次封装，使用起来非常简单高效，能够节省大量的开发时间，提高工作效率，避免重复代码。


## 具有能力
- MySql 单机或集群
- PostgreSQL 单机或集群
- SQLite 单机或集群
- SQLServer 单机或集群
- TiDB 单机或集群
- Clickhouse 单机或集群
- Redis 单机或集群
- captcha 验证码生成
- cmdc 外部命令调用
- datas 数据结构转换
- emails 邮件处理
- elasticsearch 多版本兼容封装
- files 文件工具方法封装
- gins web框架高度封装
- hook 动态代码执行
- httpc 网络请求客户端封装
- injects 依赖注入和映射
- jwts 网络请求鉴权
- kafkas 单机或集群
- machine 机器码生成
- mongoc 单机或集群
- scripts 多功能脚本
- tests 自动化测试
- times 时间日期方法封装
- utils 独立工具库
- uuid 唯一ID生成器
- 


其他工具库正在逐步完善和封装。



## 举个例子

mysql集群数据库连接使用：

```

package main

import (
	"encoding/json"
	"fmt"
	"github.com/AbnerEarl/goutils/dbs"
	"time"
)



func mai() {
	//连接数据库集群
	dsns := []string{
		"root:password@tcp(101.152.77.55:3306)/test?charset=utf8&parseTime=true",
		"root:password@tcp(101.159.16.231:3306)/test?charset=utf8&parseTime=true",
		"root:password@tcp(101.159.16.88:3306)/test?charset=utf8&parseTime=true",
	}
	db := dbs.OpenDBMySQLCluster(dsns, false, 0, 0)

	//自动创建表或者更新表结构
	db.Migration([]interface{}{&dbs.UpdateModel{}})

	//定义一个实体
	up := dbs.UpdateModel{
		FileName:    "test",
		ExecuteTime: time.Now(),
	}

	//保存数据到数据库
	err := db.Create(&up)
	fmt.Println("create result: ", err)

	//查询一条数据
	up2 := dbs.UpdateModel{}
	err = db.RetrieveByFind(&up2)
	fmt.Println("retrieve result: ", err)
	fmt.Println(up2)

	//删除数据
	err = db.DeleteHardById(&up2)
	fmt.Println("delete result: ", err)

	//批量插入数据库
	var ups = []dbs.UpdateModel{{FileName: "uu1", ExecuteTime: time.Now()}, {FileName: "uu2", ExecuteTime: time.Now()}, {FileName: "uu3", ExecuteTime: time.Now()}}
	err = db.CreateBatch(ups, 100)
	fmt.Println("batch result: ", err)

	//批量查询数据库
	bys, count, err := db.RetrieveByWhereBytes(10, 1, &dbs.UpdateModel{}, "", "", nil)
	if err != nil {
		return
	}
	fmt.Println("total count: ", count)
	var result []dbs.UpdateModel
	json.Unmarshal(bys, &result)
	fmt.Println(result)
}



```


```

go run main.go 


```



```

create result:  <nil>
retrieve result:  <nil>
{{3 0 2024-01-09 07:25:16.067 +0000 UTC 2024-01-09 07:25:16.067 +0000 UTC <nil> } test 2024-01-09 07:25:16.036 +0000 UTC}
delete result:  <nil>
batch result: <nil>
total count:  16
[{{6 0 2024-01-09 07:25:16.224 +0000 UTC 2024-01-09 07:25:16.224 +0000 UTC <nil> } uu1 2024-01-09 07:25:16.194 +0000 UTC} {{9 0 2024-01-09 07:25:16.224 +0000 UTC 2024-01-09 07:25:16.224 +0000 UTC <nil> } uu2 2024-01-09 07:25:16.194 +0000 UTC} {{12 0 2024-01-09 07:25:16.224 +0000 UTC 2024-01-09 07:25:16.224 +0000 UTC <nil> } uu3 2024-01-09 07:25:16.194 +0000 UTC} {{14 0 2024-01-09 07:35:35.615 +0000 UTC 2024-01-09 07:35:35.615 +0000 UTC <nil> } test222 2024-01-09 07:35:35.585 +0000 UTC} {{16 0 2024-01-09 07:35:35.781 +0000 UTC 2024-01-09 07:35:35.781 +0000 UTC <nil> } uu1 2024-01-09 07:35:35.751 +0000 UTC} {{18 0 2024-01-09 07:35:35.781 +0000 UTC 2024-01-09 07:35:35.781 +0000 UTC <nil> } uu2 2024-01-09 07:35:35.751 +0000 UTC} {{20 0 2024-01-09 07:35:35.781 +0000 UTC 2024-01-09 07:35:35.781 +0000 UTC <nil> } uu3 2024-01-09 07:35:35.751 +0000 UTC} {{21 0 2024-01-09 07:41:35.85 +0000 UTC 2024-01-09 07:41:35.85 +0000 UTC <nil> } test333 2024-01-09 07:41:35.819 +0000 UTC} {{23 0 2024-01-09 08:03:15.64 +0000 UTC 2024-01-09 08:03:15.64 +0000 UTC <nil> } test4444 2024-01-09 08:03:15.587 +0000 UTC} {{25 0 2024-01-09 08:03:15.881 +0000 UTC 2024-01-09 08:03:15.881 +0000 UTC <nil> } uu1 2024-01-09 08:03:15.837 +0000 UTC}]


```

## 主要方法

每个工具库都封装了若干方法，可以到对应的包下面查看，以dbs数据库使用举例，增删改查（CRUD）方法都是完备的，部分方法示例如下：

```

func (m *BaseModel) BeforeCreate(fc func(tx *TX) error, db *DB) error {
func (m *BaseModel) AfterCreate(fc func(tx *TX) error, db *DB) error {
func (m *BaseModel) BeforeSave(fc func(tx *TX) error, db *DB) error {
func (m *BaseModel) AfterSave(fc func(tx *TX) error, db *DB) error {
func (m *BaseModel) BeforeUpdate(fc func(tx *TX) error, db *DB) error {
func (m *BaseModel) AfterUpdate(fc func(tx *TX) error, db *DB) error {
func (m *BaseModel) BeforeDelete(fc func(tx *TX) error, db *DB) error {
func (m *BaseModel) AfterDelete(fc func(tx *TX) error, db *DB) error {
func (db *DB) Create(dataModel interface{}) error {
func (db *DB) CreateBatch(dataModels interface{}, batchSize uint) error {
func (db *DB) UpdateById(dataModel interface{}) error {
func (db *DB) UpdateByWhereModel(where string, updateModel interface{}) error {
func (db *DB) UpdateByArgsWhereModel(where string, args []interface{}, updateModel interface{}) error {
func (db *DB) UpdateByWhere(dataModel interface{}, where string, updates map[string]interface{}) error {
func (db *DB) UpdateByModelWhere(whereModel interface{}, updates map[string]interface{}) error {
func (db *DB) UpdateByModelWhereModel(whereModel interface{}, updateModel interface{}) error {
func (db *DB) UpdateByArgsWhere(dataModel interface{}, where string, args []interface{}, updates map[string]interface{}) error {
func (db *DB) UpdateByField(dataModel interface{}, where interface{}, column, expr string, updates ...interface{}) error {
func (db *DB) UpdateByArgsField(dataModel interface{}, where string, args []interface{}, column, expr string, updates ...interface{}) error {
func (db *DB) UpdateByModelField(whereModel interface{}, column, expr string, updates ...interface{}) error {
func (db *DB) DeleteHardById(dataModels interface{}) error {
func (db *DB) DeleteSoftById(dataModels interface{}) error {
func (db *DB) DeleteHardByWhere(dataModel interface{}, where string, args []interface{}) error {
func (db *DB) DeleteSoftByWhere(dataModel interface{}, where string, args []interface{}) error {
func (db *DB) RetrieveById(whereModel interface{}) error {
func (db *DB) RetrieveByFind(whereModel interface{}) error {
func (db *DB) RetrieveByMap(dataModel interface{}, whereMap map[string]interface{}) error {
func (db *DB) RetrieveByArgs(dataModel interface{}, where string, args []interface{}) error {
func (db *DB) RetrieveCountByArgs(dataModel interface{}, where string, args []interface{}) (int64, error) {
func (db *DB) RetrieveCountByModel(whereModel interface{}) (int64, error) {
func (db *DB) RetrieveByWhere(pageSize, pageNo int, dataModel interface{}, order, where string, args []interface{}) (interface{}, int64, error) {
func (db *DB) RetrieveByWhereString(pageSize, pageNo int, dataModel interface{}, order, where string, args []interface{}) (string, int64, error) {
func (db *DB) RetrieveByWhereBytes(pageSize, pageNo int, dataModel interface{}, order, where string, args []interface{}) ([]byte, int64, error) {
func (db *DB) RetrieveByWhereSelect(pageSize, pageNo int, dataModel interface{}, fields []string, order, where string, args []interface{}) (interface{}, int64, error) {
func (db *DB) RetrieveByWhereSelectString(pageSize, pageNo int, dataModel interface{}, fields []string, order, where string, args []interface{}) (string, int64, error) {
func (db *DB) RetrieveByWhereSelectBytes(pageSize, pageNo int, dataModel interface{}, fields []string, order, where string, args []interface{}) ([]byte, int64, error) {
func (db *DB) RawSqlForMap(sql string, values ...interface{}) ([]map[string]interface{}, error) {
func (db *DB) RawSqlForByte(sql string, values ...interface{}) ([]byte, error) {
func Rows2Map(rows *sql.Rows) []map[string]interface{} {
func Rows2Bytes(rows *sql.Rows) []byte {
func (db *DB) Exec(sql string, values ...interface{}) error {
func (db *DB) RetrieveByModel(pageSize, pageNo int, whereModel interface{}, order string) (interface{}, int64, error) {
func (db *DB) RetrieveByModelString(pageSize, pageNo int, whereModel interface{}, order string) (string, int64, error) {
func (db *DB) RetrieveByModelBytes(pageSize, pageNo int, whereModel interface{}, order string) ([]byte, int64, error) {
func (db *DB) RetrieveByModelSelect(pageSize, pageNo int, whereModel interface{}, fields []string, order string) (interface{}, int64, error) {
func (db *DB) RetrieveByModelSelectString(pageSize, pageNo int, whereModel interface{}, fields []string, order string) (string, int64, error) {
func (db *DB) RetrieveByModelSelectBytes(pageSize, pageNo int, whereModel interface{}, fields []string, order string) ([]byte, int64, error) {
func (db *DB) Transaction(fc func(tx *TX) error) error {


```


复杂的工具库使用方法都可以到 docs 目录下查看相关文档。



## 相关说明

欢迎有兴趣的朋友一起参与，或提出宝贵建议。



