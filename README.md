
## [English](https://github.com/AbnerEarl/goutils/blob/main/README.md) ｜  [中文](https://github.com/AbnerEarl/goutils/blob/main/Chinese.md)

# goutils
Some very useful middleware libraries are re-encapsulated and are very simple and efficient to use, which can save a lot of development time, improve work efficiency and avoid duplication of code.


## Have the ability
- MySql stand-alone or cluster
- PostgreSQL standalone or cluster
- SQLite standalone or cluster
- SQLServer stand-alone or cluster
- TiDB stand-alone or cluster
- Clickhouse standalone or cluster
- Redis stand-alone or cluster
- captcha verification code generation
- cmdc external command call
- datas data structure conversion
- emails mail processing
- elasticsearch multi-version compatible package
- files file tool method encapsulation
- The gins web framework is highly encapsulated
- hook dynamic code execution
- httpc network request client encapsulation
- injects dependency injection and mapping
- jwts network request authentication
  -kafkas standalone or cluster
- machine machine code generation
- mongoc stand-alone or cluster
- scripts multifunctional scripts
- tests automated tests
- times time and date method encapsulation
- utils independent tool library
- uuid unique ID generator
-


Other tool libraries are being gradually improved and encapsulated.




















































## for example

Mysql cluster database connection uses:

```

package main

import (
"encoding/json"
"fmt"
"github.com/AbnerEarl/goutils/dbs"
"time"
)



func mai() {
//Connect to the database cluster
dsns := []string{
"root:password@tcp(101.152.77.55:3306)/test?charset=utf8&parseTime=true",
"root:password@tcp(101.159.16.231:3306)/test?charset=utf8&parseTime=true",
"root:password@tcp(101.159.16.88:3306)/test?charset=utf8&parseTime=true",
}
db := dbs.OpenDBMySQLCluster(dsns, false, 0, 0)

//Automatically create a table or update the table structure
db.Migration([]interface{}{&dbs.UpdateModel{}})

//Define an entity
up := dbs.UpdateModel{
FileName: "test",
ExecuteTime: time.Now(),
}

//Save data to database
err := db.Create(&up)
fmt.Println("create result: ", err)

//Query a piece of data
up2 := dbs.UpdateModel{}
err = db.RetrieveByFind(&up2)
fmt.Println("retrieve result: ", err)
fmt.Println(up2)

//delete data
err = db.DeleteHardById(&up2)
fmt.Println("delete result: ", err)

//Batch insert into database
var ups = []dbs.UpdateModel{{FileName: "uu1", ExecuteTime: time.Now()}, {FileName: "uu2", ExecuteTime: time.Now()}, {FileName: "uu3", ExecuteTime: time.Now()}}
err = db.CreateBatch(ups, 100)
fmt.Println("batch result: ", err)

//Batch query database
bys, count, err := db.RetrieveByWhereBytes(10, 1, &dbs.UpdateModel{}, "", "", nil)
if err != nil {
return
}
fmt.Println("total count: ", count)
var result[]dbs.UpdateModel
json.Unmarshal(bys, &result)
fmt.Println(result)
}



```


```

go run main.go


```



```

create result: <nil>
retrieve result: <nil>
{{3 0 2024-01-09 07:25:16.067 +0000 UTC 2024-01-09 07:25:16.067 +0000 UTC <nil> } test 2024-01-09 07:25:16.036 +0000 UTC}
delete result: <nil>
batch result: <nil>
total count: 16
[{{6 0 2024-01-09 07:25:16.224 +0000 UTC 2024-01-09 07:25:16.224 +0000 UTC <nil> } uu1 2024-01-09 07:25:16.194 +0000 UTC} {{9 0 2024-01-09 07:25:16.224 +0000 UTC 2024-01-09 07:25:16.224 +0000 UTC <nil> } uu2 2024-01-09 07:25:16.194 +0000 UTC} { {12 0 2024-01-09 07:25:16.224 +0000 UTC 2024-01-09 07:25:16.224 +0000 UTC <nil> } uu3 2024-01-09 07:25:16.194 +0000 UTC} {{ 14 0 2024-01-09 07:35:35.615 +0000 UTC 2024-01-09 07:35:35.615 +0000 UTC <nil> } test222 2024-01-09 07:35:35.585 +0000 UTC} {{16 0 2024-01-09 07:35:35.781 +0000 UTC 2024-01-09 07:35:35.781 +0000 UTC <nil> } uu1 2024-01-09 07:35:35.751 +0000 UTC} {{18 0 2024-01-09 07:35:35.781 +0000 UTC 2024-01-09 07:35:35.781 +0000 UTC <nil> } uu2 2024-01-09 07:35:35.751 +0000 UTC} {{20 0 2024 -01-09 07:35:35.781 +0000 UTC 2024-01-09 07:35:35.781 +0000 UTC <nil> } uu3 2024-01-09 07:35:35.751 +0000 UTC} {{21 0 2024- 01-09 07:41:35.85 +0000 UTC 2024-01-09 07:41:35.85 +0000 UTC <nil> } test333 2024-01-09 07:41:35.819 +0000 UTC} {{23 0 2024-01 -09 08:03:15.64 +0000 UTC 2024-01-09 08:03:15.64 +0000 UTC <nil> } test4444 2024-01-09 08:03:15.587 +0000 UTC} {{25 0 2024-01- 09 08:03:15.881 +0000 UTC 2024-01-09 08:03:15.881 +0000 UTC <nil> } uu1 2024-01-09 08:03:15.837 +0000 UTC}]


```



## Main method

Each tool library encapsulates a number of methods, which can be viewed under the corresponding package. Taking the use of dbs database as an example, the Create,Retrieve,Update,Delete (CRUD) method is complete. Some method examples are as follows:

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


For complex tool library usage, you can view relevant documents in the docs directory.



## Related instructions

Interested friends are welcome to participate or make valuable suggestions.

