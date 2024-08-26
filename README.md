
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



In order to better use automatic association methods, it is necessary to pay attention to some rules, such as first defining a model:

```
/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/8/14 15:55
 * @desc: about the role of class.
 */

package model

import (
	"freedom/config"
	"github.com/AbnerEarl/goutils/dbs"
	"github.com/AbnerEarl/goutils/web"
)

type ShopModel struct {
	dbs.BaseModel
	UserId        uint64  `json:"user_id" gorm:"column:user_id;not null;comment:'用户ID'"`
	Name          string  `json:"name" gorm:"column:name;not null;comment:'名称'"`
	Image         string  `json:"image" gorm:"column:image;not null;comment:'店铺照片'"`
	Location      string  `json:"location" gorm:"column:location;not null;comment:'位置'"`
	Longitude     float64 `json:"longitude" gorm:"column:longitude;not null;comment:'经度'"`
	Latitude      float64 `json:"latitude" gorm:"column:latitude;not null;comment:'纬度'"`
	Phone         string  `json:"phone" gorm:"column:phone;not null;comment:'电话'"`
	WorkTime      string  `json:"work_time" gorm:"column:work_time;not null;comment:'营业时间'"`
	WorkStartTime string  `json:"work_start_time" gorm:"column:work_start_time;not null;comment:'营业开始时间'"`
	WorkEndTime   string  `json:"work_end_time" gorm:"column:work_end_time;not null;comment:'营业结束时间'"`
	WorkStatus    uint    `json:"work_status" gorm:"column:work_status;default:0;not null;comment:'营业状态'"`
	Status        uint    `json:"status" gorm:"column:status;default:1;not null;comment:'店铺状态'"`
}

func (m *ShopModel) TableName() string {
	return "shop_info"
}

// For the sake of data security and efficiency, an irregular monotonically increasing integer is used as the primary key ID,
// which is mainly improved based on the Snow Flower algorithm, Pear Blossom algorithm, and Mist algorithm.
func (m *ShopModel) BeforeCreate() error {
	m.Id = web.GenAutoIdByKeyClu(config.RedisCli, web.AUTO_ID_GENERATOR_COUNTER_KEY+m.TableName())
	return nil
}

```


There is no need to define other methods, as this model already inherits most of the CRUD related methods.

Automatically generate relevant documents, including dictionaries, API documents, etc. You can use them by creating a gen.go file in the root directory or main file directory, and entering the following content:

```
/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/7/4 09:21
 * @desc: about the role of class.
 */

package main

import (
	"fmt"
	"github.com/AbnerEarl/goutils/files"
	"github.com/AbnerEarl/goutils/scripts"
	"github.com/AbnerEarl/goutils/swag"
)

func main() {
	err := scripts.GenDbComment(files.GetCuPath()+"/model", "base", "comment.go")
	if err != nil {
		fmt.Println("GenDbComment: ", err)
	}

	err = swag.InstallSwag()
	if err != nil {
		fmt.Println("InstallSwag: ", err)
	}

	err = swag.GenSwagDoc("shop.go", "./api/shop", nil)
	if err != nil {
		fmt.Println("GenSwagDoc: ", err)
	}

}

```


Finally, execute the gen.go file to generate the relevant dictionary and API interface documents, as shown in the following example:

```
go run gen.go 

```

----

For other tool library usage, you can view relevant documents in the docs directory.



## Related instructions

Interested friends are welcome to participate or make valuable suggestions.

