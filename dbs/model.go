package dbs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/AbnerEarl/goutils/utils"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

var (
	DefaultLimit    = 10
	DefaultMaxLimit = 1000
)

type TX struct {
	*gorm.DB
}

type BaseModel struct {
	Id        uint64     `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id;comment:'主键ID'"`
	IsDel     uint64     `json:"-" gorm:"column:is_del;default:0;comment:'删除标志'"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;comment:'创建时间'"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at;comment:'更新时间'"`
	DeletedAt *time.Time `json:"-" gorm:"column:deleted_at;comment:'删除时间'" sql:"index"`
	Remark    string     `json:"remark" gorm:"column:remark;null;type:text;comment:'备注信息'"`
}

func (m *BaseModel) TableName() string {
	return "base_info"
}

func (m *BaseModel) BeforeCreate(fc func(tx *TX) error, db *DB) error {
	f := HandleFunc(fc)
	return db.DB.Transaction(f)
}

func (m *BaseModel) AfterCreate(fc func(tx *TX) error, db *DB) error {
	f := HandleFunc(fc)
	return db.DB.Transaction(f)
}

func (m *BaseModel) BeforeSave(fc func(tx *TX) error, db *DB) error {
	f := HandleFunc(fc)
	return db.DB.Transaction(f)
}

func (m *BaseModel) AfterSave(fc func(tx *TX) error, db *DB) error {
	f := HandleFunc(fc)
	return db.DB.Transaction(f)
}

func (m *BaseModel) BeforeUpdate(fc func(tx *TX) error, db *DB) error {
	f := HandleFunc(fc)
	return db.DB.Transaction(f)
}

func (m *BaseModel) AfterUpdate(fc func(tx *TX) error, db *DB) error {
	f := HandleFunc(fc)
	return db.DB.Transaction(f)
}

func (m *BaseModel) BeforeDelete(fc func(tx *TX) error, db *DB) error {
	f := HandleFunc(fc)
	return db.DB.Transaction(f)
}

func (m *BaseModel) AfterDelete(fc func(tx *TX) error, db *DB) error {
	f := HandleFunc(fc)
	return db.DB.Transaction(f)
}

func (db *DB) Create(dataModel interface{}) error {
	//tableName := ""
	//ref := reflect.ValueOf(dataModel)
	//method := ref.MethodByName("TableName")
	//if method.IsValid() {
	//	r := method.Call([]reflect.Value{})
	//	tableName = r[0].String()
	//} else {
	//	return fmt.Errorf("the current model does not have a table name defined")
	//}
	//return DB.Table(tableName).Create(dataModel).Error
	return db.DB.Create(dataModel).Error
}

func (db *DB) CreateBatch(dataModels interface{}, batchSize uint) error {
	return db.CreateInBatches(dataModels, int(batchSize)).Error
}

func (db *DB) UpdateById(dataModel interface{}) error {
	return db.Save(dataModel).Error
}

func (db *DB) UpdateByWhereModel(where string, updateModel interface{}) error {
	tableName := ""
	ref := reflect.ValueOf(updateModel)
	method := ref.MethodByName("TableName")
	if method.IsValid() {
		r := method.Call([]reflect.Value{})
		tableName = r[0].String()
	} else {
		return fmt.Errorf("the current model does not have a table name defined")
	}
	return db.Table(tableName).Where(where).Updates(updateModel).Error
}

func (db *DB) UpdateByArgsWhereModel(where string, args []interface{}, updateModel interface{}) error {
	tableName := ""
	ref := reflect.ValueOf(updateModel)
	method := ref.MethodByName("TableName")
	if method.IsValid() {
		r := method.Call([]reflect.Value{})
		tableName = r[0].String()
	} else {
		return fmt.Errorf("the current model does not have a table name defined")
	}
	return db.Table(tableName).Where(where, args...).Updates(updateModel).Error
}

func (db *DB) UpdateByWhere(dataModel interface{}, where string, updates map[string]interface{}) error {
	return db.Model(dataModel).Where(where).Updates(updates).Error
}

func (db *DB) UpdateByModelWhere(whereModel interface{}, updates map[string]interface{}) error {
	return db.Model(whereModel).Where(whereModel).Updates(updates).Error
}

func (db *DB) UpdateByModelWhereModel(whereModel interface{}, updateModel interface{}) error {
	return db.Model(whereModel).Where(whereModel).Updates(updateModel).Error
}

func (db *DB) UpdateByArgsWhere(dataModel interface{}, where string, args []interface{}, updates map[string]interface{}) error {
	return db.Model(dataModel).Where(where, args...).Updates(updates).Error
}

func (db *DB) UpdateByField(dataModel interface{}, where interface{}, column, expr string, updates ...interface{}) error {
	//use example:
	//m := UpdateModel{}
	//err := UpdateByArgsField(&m, "id = ?", 1, "value", "value + ?", 1)
	return db.Model(dataModel).Where(where).Update(column, gorm.Expr(expr, updates...)).Error
}

func (db *DB) UpdateByArgsField(dataModel interface{}, where string, args []interface{}, column, expr string, updates ...interface{}) error {
	//use example:
	//m := UpdateModel{}
	//err := UpdateByArgsField(&m, "id = ?", []interface{}{1}, "value", "value + ?", 1)
	return db.Model(dataModel).Where(where, args...).Update(column, gorm.Expr(expr, updates...)).Error
}

func (db *DB) UpdateByModelField(whereModel interface{}, column, expr string, updates ...interface{}) error {
	//use example:
	//m := UpdateModel{}
	//m.Id = 1
	//err := UpdateByModelField(&m, "value", "value + ?", 1)
	return db.Model(whereModel).Where(whereModel).Update(column, gorm.Expr(expr, updates...)).Error
}

func (db *DB) DeleteHardById(dataModels interface{}) error {
	return db.Unscoped().Delete(dataModels).Error
}

func (db *DB) DeleteSoftById(dataModels interface{}) error {
	return db.Delete(dataModels).Error
}

func (db *DB) DeleteHardByWhere(dataModel interface{}, where string, args []interface{}) error {
	return db.Unscoped().Where(where, args...).Delete(dataModel).Error
}

func (db *DB) DeleteSoftByWhere(dataModel interface{}, where string, args []interface{}) error {
	return db.Where(where, args...).Delete(dataModel).Error
}

func (db *DB) RetrieveFirstById(whereModel interface{}) error {
	return db.First(whereModel).Error
}

func (db *DB) RetrieveFirstByFind(whereModel interface{}) error {
	return db.Where(whereModel).First(whereModel).Error
}

func (db *DB) RetrieveFirstByMap(dataModel interface{}, whereMap map[string]interface{}) error {
	return db.Where(whereMap).First(dataModel).Error
}

func (db *DB) RetrieveFirstByArgs(dataModel interface{}, where string, args []interface{}) error {
	return db.Where(where, args...).First(dataModel).Error
}

func (db *DB) RetrieveLastById(whereModel interface{}) error {
	return db.Last(whereModel).Error
}

func (db *DB) RetrieveLastByFind(whereModel interface{}) error {
	return db.Where(whereModel).Last(whereModel).Error
}

func (db *DB) RetrieveLastByMap(dataModel interface{}, whereMap map[string]interface{}) error {
	return db.Where(whereMap).Last(dataModel).Error
}

func (db *DB) RetrieveLastByArgs(dataModel interface{}, where string, args []interface{}) error {
	return db.Where(where, args...).Last(dataModel).Error
}

func (db *DB) RetrieveCountByArgs(dataModel interface{}, where string, args []interface{}) (int64, error) {
	var count int64 = 0
	err := db.Model(dataModel).Where(where, args...).Count(&count).Error
	return count, err
}

func (db *DB) RetrieveCountByModel(whereModel interface{}) (int64, error) {
	var count int64 = 0
	err := db.Model(whereModel).Where(whereModel).Count(&count).Error
	return count, err
}

func (db *DB) RetrieveByWhere(pageSize, pageNo int, dataModel interface{}, order, where string, args []interface{}) (interface{}, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveByWhere(0, 0, &m, "", "id=?", []interface{}{1})
	//dataList := result.(*[]*dbs.UpdateModel)
	//fmt.Println((*data)[0])
	var count int64 = 0
	if err := db.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(dataModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	if err := db.
		Model(dataModel).
		Where(where, args...).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	return results.Interface(), count, nil
}

func (db *DB) RetrieveOrderInByWhere(pageSize, pageNo int, dataModel interface{}, where string, args []interface{}, orderKey string, orderValues []interface{}) (interface{}, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveByWhere(0, 0, &m, "id=?", []interface{}{1}, "id", []interface{}{3,2,1})
	//dataList := result.(*[]*dbs.UpdateModel)
	//fmt.Println((*data)[0])
	var count int64 = 0
	if err := db.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(dataModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	order := ""
	if len(orderValues) > 0 {
		order = fmt.Sprintf("FIELD(%s,%s)", orderKey, utils.List2String(orderValues))
	}
	if err := db.
		Model(dataModel).
		Where(where, args...).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	return results.Interface(), count, nil
}

func (db *DB) RetrieveByWhereString(pageSize, pageNo int, dataModel interface{}, order, where string, args []interface{}) (string, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveByWhereString(0, 0, &m, "", "id=?", []interface{}{1})
	//fmt.Println(result)
	var count int64 = 0
	if err := db.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
		return "", count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(dataModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	if err := db.
		Model(dataModel).
		Where(where, args...).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return "", count, err
	}
	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return "", count, err
	}
	return string(bytes), count, nil
}

func (db *DB) RetrieveOrderInByWhereString(pageSize, pageNo int, dataModel interface{}, where string, args []interface{}, orderKey string, orderValues []interface{}) (string, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveByWhereString(0, 0, &m, "id=?", []interface{}{1}, "id", []interface{}{3,2,1})
	//fmt.Println(result)
	var count int64 = 0
	if err := db.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
		return "", count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(dataModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	order := ""
	if len(orderValues) > 0 {
		order = fmt.Sprintf("FIELD(%s,%s)", orderKey, utils.List2String(orderValues))
	}
	if err := db.
		Model(dataModel).
		Where(where, args...).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return "", count, err
	}
	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return "", count, err
	}
	return string(bytes), count, nil
}

func (db *DB) RetrieveByWhereBytes(pageSize, pageNo int, dataModel interface{}, order, where string, args []interface{}) ([]byte, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveByWhereBytes(0, 0, &m, "", "id=?", []interface{}{1})
	//var dataList []interface{}
	//json.Unmarshal(result,&dataList)
	//fmt.Println(dataList)
	var count int64 = 0
	if err := db.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(dataModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	if err := db.
		Model(dataModel).
		Where(where, args...).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return nil, count, err
	}
	return bytes, count, nil
}

func (db *DB) RetrieveOrderInByWhereBytes(pageSize, pageNo int, dataModel interface{}, where string, args []interface{}, orderKey string, orderValues []interface{}) ([]byte, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveByWhereBytes(0, 0, &m, "id=?", []interface{}{1}, "id", []interface{}{3,2,1})
	//var dataList []interface{}
	//json.Unmarshal(result,&dataList)
	//fmt.Println(dataList)
	var count int64 = 0
	if err := db.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(dataModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	order := ""
	if len(orderValues) > 0 {
		order = fmt.Sprintf("FIELD(%s,%s)", orderKey, utils.List2String(orderValues))
	}
	if err := db.
		Model(dataModel).
		Where(where, args...).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return nil, count, err
	}
	return bytes, count, nil
}
func (db *DB) RetrieveByWhereSelect(pageSize, pageNo int, dataModel interface{}, fields []string, order, where string, args []interface{}) (interface{}, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveBySelect(0, 0, &m, []string{"id", "name"}, "", "id=?", []interface{}{1})
	//dataList := result.(*[]*dbs.UpdateModel)
	//fmt.Println((*dataList)[0])
	var count int64
	if err := db.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(dataModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	if err := db.
		Select(fields).
		Where(where, args...).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	return results.Interface(), count, nil
}

func (db *DB) RetrieveOrderInByWhereSelect(pageSize, pageNo int, dataModel interface{}, fields []string, where string, args []interface{}, orderKey string, orderValues []interface{}) (interface{}, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveBySelect(0, 0, &m, []string{"id", "name"}, "id=?", []interface{}{1}, "id", []interface{}{3,2,1})
	//dataList := result.(*[]*dbs.UpdateModel)
	//fmt.Println((*dataList)[0])
	var count int64
	if err := db.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(dataModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	order := ""
	if len(orderValues) > 0 {
		order = fmt.Sprintf("FIELD(%s,%s)", orderKey, utils.List2String(orderValues))
	}
	if err := db.
		Select(fields).
		Where(where, args...).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	return results.Interface(), count, nil
}
func (db *DB) RetrieveByWhereSelectString(pageSize, pageNo int, dataModel interface{}, fields []string, order, where string, args []interface{}) (string, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveBySelectString(0, 0, &m, []string{"id", "name"}, "", "id=?", []interface{}{1})
	//fmt.Println(result)
	var count int64
	if err := db.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
		return "", count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(dataModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	if err := db.
		Select(fields).
		Where(where, args...).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return "", count, err
	}

	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return "", count, err
	}
	return string(bytes), count, nil
}

func (db *DB) RetrieveOrderInByWhereSelectString(pageSize, pageNo int, dataModel interface{}, fields []string, where string, args []interface{}, orderKey string, orderValues []interface{}) (string, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveBySelectString(0, 0, &m, []string{"id", "name"}, "id=?", []interface{}{1}, "id", []interface{}{3,2,1})
	//fmt.Println(result)
	var count int64
	if err := db.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
		return "", count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(dataModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	order := ""
	if len(orderValues) > 0 {
		order = fmt.Sprintf("FIELD(%s,%s)", orderKey, utils.List2String(orderValues))
	}
	if err := db.
		Select(fields).
		Where(where, args...).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return "", count, err
	}

	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return "", count, err
	}
	return string(bytes), count, nil
}
func (db *DB) RetrieveByWhereSelectBytes(pageSize, pageNo int, dataModel interface{}, fields []string, order, where string, args []interface{}) ([]byte, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveBySelectBytes(0, 0, &m, []string{"id", "name"}, "", "id=?", []interface{}{1})
	//var dataList []interface{}
	//json.Unmarshal(result,&dataList)
	//fmt.Println(dataList)
	var count int64
	if err := db.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(dataModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	if err := db.
		Select(fields).
		Where(where, args...).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}

	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return nil, count, err
	}
	return bytes, count, nil
}

func (db *DB) RetrieveOrderInByWhereSelectBytes(pageSize, pageNo int, dataModel interface{}, fields []string, where string, args []interface{}, orderKey string, orderValues []interface{}) ([]byte, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveBySelectBytes(0, 0, &m, []string{"id", "name"}, "id=?", []interface{}{1}, "id", []interface{}{3,2,1})
	//var dataList []interface{}
	//json.Unmarshal(result,&dataList)
	//fmt.Println(dataList)
	var count int64
	if err := db.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(dataModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	order := ""
	if len(orderValues) > 0 {
		order = fmt.Sprintf("FIELD(%s,%s)", orderKey, utils.List2String(orderValues))
	}
	if err := db.
		Select(fields).
		Where(where, args...).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}

	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return nil, count, err
	}
	return bytes, count, nil
}

func (db *DB) RawSqlForMap(sql string, values ...interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	rows, err := db.Raw(sql, values...).Rows()
	if err != nil {
		return result, err
	}
	result = Rows2Map(rows)
	return result, nil
}

func (db *DB) RawSqlForByte(sql string, values ...interface{}) ([]byte, error) {
	rows, err := db.Raw(sql, values...).Rows()
	if err != nil {
		return nil, err
	}
	result := Rows2Bytes(rows)
	return result, nil
}

func Rows2Map(rows *sql.Rows) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	colTypes, _ := rows.ColumnTypes()
	var rowParam = make([]interface{}, len(colTypes))
	var rowValue = make([]interface{}, len(colTypes))

	for i, colType := range colTypes {
		rowValue[i] = reflect.New(colType.ScanType())
		rowParam[i] = reflect.ValueOf(&rowValue[i]).Interface()
	}

	for rows.Next() {
		rows.Scan(rowParam...)
		record := make(map[string]interface{})
		for i, colType := range colTypes {
			if rowValue[i] == nil {
				record[colType.Name()] = rowValue[i]
			} else {
				if strings.Contains(reflect.TypeOf(rowValue[i]).String(), "[]uint8") {
					record[colType.Name()] = utils.Byte2Any(rowValue[i].([]byte), colType.ScanType())
				} else {
					val, err := json.Marshal(rowValue[i])
					if err != nil {
						record[colType.Name()] = rowValue[i]
					} else {
						record[colType.Name()] = utils.Byte2Any(val, colType.ScanType())
					}
				}
			}
		}
		res = append(res, record)
	}
	return res
}

func Rows2Bytes(rows *sql.Rows) []byte {
	res := make([]map[string]interface{}, 0)
	colTypes, _ := rows.ColumnTypes()
	var rowParam = make([]interface{}, len(colTypes))
	var rowValue = make([]interface{}, len(colTypes))

	for i, colType := range colTypes {
		rowValue[i] = reflect.New(colType.ScanType())
		rowParam[i] = reflect.ValueOf(&rowValue[i]).Interface()
	}

	for rows.Next() {
		rows.Scan(rowParam...)
		record := make(map[string]interface{})
		for i, colType := range colTypes {
			if rowValue[i] != nil && strings.Contains(reflect.TypeOf(rowValue[i]).String(), "[]uint8") {
				record[colType.Name()] = utils.Byte2Any(rowValue[i].([]byte), colType.ScanType())
			} else {
				record[colType.Name()] = rowValue[i]
			}
		}
		res = append(res, record)
	}
	result, _ := json.Marshal(res)
	return result
}

func (db *DB) Exec(sql string, values ...interface{}) error {
	return db.DB.Exec(sql, values...).Error
}

func (db *DB) RetrieveByModel(pageSize, pageNo int, whereModel interface{}, order string) (interface{}, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModel(0, 0, &m, "")
	//dataList := result.(*[]*dbs.UpdateModel)
	//fmt.Println((*data)[0])
	var count int64 = 0
	if err := db.Model(whereModel).Where(whereModel).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(whereModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	if err := db.
		Model(whereModel).
		Where(whereModel).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	return results.Interface(), count, nil
}

func (db *DB) RetrieveOrderInByModel(pageSize, pageNo int, whereModel interface{}, orderKey string, orderValues []interface{}) (interface{}, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModel(0, 0, &m, "id", []interface{}{3,2,1})
	//dataList := result.(*[]*dbs.UpdateModel)
	//fmt.Println((*data)[0])
	var count int64 = 0
	if err := db.Model(whereModel).Where(whereModel).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(whereModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	order := ""
	if len(orderValues) > 0 {
		order = fmt.Sprintf("FIELD(%s,%s)", orderKey, utils.List2String(orderValues))
	}
	if err := db.
		Model(whereModel).
		Where(whereModel).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	return results.Interface(), count, nil
}

func (db *DB) RetrieveByModelString(pageSize, pageNo int, whereModel interface{}, order string) (string, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelString(0, 0, &m, "")
	//fmt.Println(result)
	var count int64 = 0
	if err := db.Model(whereModel).Where(whereModel).Count(&count).Error; err != nil {
		return "", count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(whereModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	if err := db.
		Model(whereModel).
		Where(whereModel).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return "", count, err
	}
	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return "", count, err
	}
	return string(bytes), count, nil
}

func (db *DB) RetrieveOrderInByModelString(pageSize, pageNo int, whereModel interface{}, orderKey string, orderValues []interface{}) (string, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelString(0, 0, &m, "id", []interface{}{3,2,1})
	//fmt.Println(result)
	var count int64 = 0
	if err := db.Model(whereModel).Where(whereModel).Count(&count).Error; err != nil {
		return "", count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(whereModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	order := ""
	if len(orderValues) > 0 {
		order = fmt.Sprintf("FIELD(%s,%s)", orderKey, utils.List2String(orderValues))
	}
	if err := db.
		Model(whereModel).
		Where(whereModel).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return "", count, err
	}
	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return "", count, err
	}
	return string(bytes), count, nil
}

func (db *DB) RetrieveByModelBytes(pageSize, pageNo int, whereModel interface{}, order string) ([]byte, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelBytes(0, 0, &m, "")
	//var dataList []interface{}
	//json.Unmarshal(result,&dataList)
	//fmt.Println(dataList)
	var count int64 = 0
	if err := db.Model(whereModel).Where(whereModel).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(whereModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	if err := db.
		Model(whereModel).
		Where(whereModel).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return nil, count, err
	}
	return bytes, count, nil
}

func (db *DB) RetrieveOrderInByModelBytes(pageSize, pageNo int, whereModel interface{}, orderKey string, orderValues []interface{}) ([]byte, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelBytes(0, 0, &m, "id", []interface{}{3,2,1})
	//var dataList []interface{}
	//json.Unmarshal(result,&dataList)
	//fmt.Println(dataList)
	var count int64 = 0
	if err := db.Model(whereModel).Where(whereModel).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(whereModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	order := ""
	if len(orderValues) > 0 {
		order = fmt.Sprintf("FIELD(%s,%s)", orderKey, utils.List2String(orderValues))
	}
	if err := db.
		Model(whereModel).
		Where(whereModel).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return nil, count, err
	}
	return bytes, count, nil
}
func (db *DB) RetrieveByModelSelect(pageSize, pageNo int, whereModel interface{}, fields []string, order string) (interface{}, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelSelect(0, 0, &m, []string{"id", "name"}, "")
	//dataList := result.(*[]*dbs.UpdateModel)
	//fmt.Println((*dataList)[0])
	var count int64
	if err := db.Model(whereModel).Where(whereModel).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(whereModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	if err := db.
		Model(whereModel).
		Where(whereModel).
		Select(fields).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	return results.Interface(), count, nil
}

func (db *DB) RetrieveOrderInByModelSelect(pageSize, pageNo int, whereModel interface{}, fields []string, orderKey string, orderValues []interface{}) (interface{}, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelSelect(0, 0, &m, []string{"id", "name"}, "id", []interface{}{3,2,1})
	//dataList := result.(*[]*dbs.UpdateModel)
	//fmt.Println((*dataList)[0])
	var count int64
	if err := db.Model(whereModel).Where(whereModel).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(whereModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	order := ""
	if len(orderValues) > 0 {
		order = fmt.Sprintf("FIELD(%s,%s)", orderKey, utils.List2String(orderValues))
	}
	if err := db.
		Model(whereModel).
		Where(whereModel).
		Select(fields).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	return results.Interface(), count, nil
}

func (db *DB) RetrieveByModelSelectString(pageSize, pageNo int, whereModel interface{}, fields []string, order string) (string, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelSelectString(0, 0, &m, []string{"id", "name"}, "")
	//fmt.Println(result)
	var count int64
	if err := db.Model(whereModel).Where(whereModel).Count(&count).Error; err != nil {
		return "", count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(whereModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	if err := db.
		Model(whereModel).
		Where(whereModel).
		Select(fields).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return "", count, err
	}

	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return "", count, err
	}
	return string(bytes), count, nil
}

func (db *DB) RetrieveOrderInByModelSelectString(pageSize, pageNo int, whereModel interface{}, fields []string, orderKey string, orderValues []interface{}) (string, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelSelectString(0, 0, &m, []string{"id", "name"}, "id", []interface{}{3,2,1})
	//fmt.Println(result)
	var count int64
	if err := db.Model(whereModel).Where(whereModel).Count(&count).Error; err != nil {
		return "", count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(whereModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	order := ""
	if len(orderValues) > 0 {
		order = fmt.Sprintf("FIELD(%s,%s)", orderKey, utils.List2String(orderValues))
	}
	if err := db.
		Model(whereModel).
		Where(whereModel).
		Select(fields).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return "", count, err
	}

	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return "", count, err
	}
	return string(bytes), count, nil
}

func (db *DB) RetrieveByModelSelectBytes(pageSize, pageNo int, whereModel interface{}, fields []string, order string) ([]byte, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelSelectBytes(0, 0, &m, []string{"id", "name"}, "")
	//var dataList []interface{}
	//json.Unmarshal(result,&dataList)
	//fmt.Println(dataList)
	var count int64
	if err := db.Model(whereModel).Where(whereModel).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(whereModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	if err := db.
		Model(whereModel).
		Where(whereModel).
		Select(fields).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}

	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return nil, count, err
	}
	return bytes, count, nil
}

func (db *DB) RetrieveOrderInByModelSelectBytes(pageSize, pageNo int, whereModel interface{}, fields []string, orderKey string, orderValues []interface{}) ([]byte, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelSelectBytes(0, 0, &m, []string{"id", "name"}, "id", []interface{}{3,2,1})
	//var dataList []interface{}
	//json.Unmarshal(result,&dataList)
	//fmt.Println(dataList)
	var count int64
	if err := db.Model(whereModel).Where(whereModel).Count(&count).Error; err != nil {
		return nil, count, err
	}
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	typ := reflect.TypeOf(whereModel)
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PtrTo(typ)
	}
	itemSlice := reflect.SliceOf(typ)
	results := reflect.New(itemSlice)
	order := ""
	if len(orderValues) > 0 {
		order = fmt.Sprintf("FIELD(%s,%s)", orderKey, utils.List2String(orderValues))
	}
	if err := db.
		Model(whereModel).
		Where(whereModel).
		Select(fields).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}

	bytes, err := json.Marshal(results.Interface())
	if err != nil {
		return nil, count, err
	}
	return bytes, count, nil
}

func (db *DB) Transaction(fc func(tx *TX) error) error {
	f := HandleFunc(fc)
	return db.DB.Transaction(f)
}

func HandleFunc(handler func(tx *TX) error) func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		return handler(&TX{db})
	}
}

func (db *DB) TruncateTable(dataModel interface{}) error {
	tableName := ""
	ref := reflect.ValueOf(dataModel)
	method := ref.MethodByName("TableName")
	if method.IsValid() {
		r := method.Call([]reflect.Value{})
		tableName = r[0].String()
	} else {
		return fmt.Errorf("the current model does not have a table name defined")
	}
	rSql := fmt.Sprintf("TRUNCATE TABLE %s;", tableName)
	err := db.Exec(rSql)
	return err
}
