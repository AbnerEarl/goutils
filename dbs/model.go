package dbs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/YouAreOnlyOne/goutils/utils"
	"gorm.io/gorm"
	"reflect"
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

func (m *BaseModel) BeforeCreate(fc func(tx *TX) error) error {
	f := HandleFunc(fc)
	return DB.Transaction(f)
}

func (m *BaseModel) AfterCreate(fc func(tx *TX) error) error {
	f := HandleFunc(fc)
	return DB.Transaction(f)
}

func (m *BaseModel) BeforeSave(fc func(tx *TX) error) error {
	f := HandleFunc(fc)
	return DB.Transaction(f)
}

func (m *BaseModel) AfterSave(fc func(tx *TX) error) error {
	f := HandleFunc(fc)
	return DB.Transaction(f)
}

func (m *BaseModel) BeforeUpdate(fc func(tx *TX) error) error {
	f := HandleFunc(fc)
	return DB.Transaction(f)
}

func (m *BaseModel) AfterUpdate(fc func(tx *TX) error) error {
	f := HandleFunc(fc)
	return DB.Transaction(f)
}

func (m *BaseModel) BeforeDelete(fc func(tx *TX) error) error {
	f := HandleFunc(fc)
	return DB.Transaction(f)
}

func (m *BaseModel) AfterDelete(fc func(tx *TX) error) error {
	f := HandleFunc(fc)
	return DB.Transaction(f)
}

func Create(dataModel interface{}) error {
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
	return DB.Create(dataModel).Error
}

func CreateBatch(dataModels interface{}, batchSize uint) error {
	return DB.CreateInBatches(dataModels, int(batchSize)).Error
}

func UpdateById(dataModel interface{}) error {
	return DB.Save(dataModel).Error
}

func UpdateByWhereModel(where string, updateModel interface{}) error {
	tableName := ""
	ref := reflect.ValueOf(updateModel)
	method := ref.MethodByName("TableName")
	if method.IsValid() {
		r := method.Call([]reflect.Value{})
		tableName = r[0].String()
	} else {
		return fmt.Errorf("the current model does not have a table name defined")
	}
	return DB.Table(tableName).Where(where).Updates(updateModel).Error
}

func UpdateByArgsWhereModel(where string, args []interface{}, updateModel interface{}) error {
	tableName := ""
	ref := reflect.ValueOf(updateModel)
	method := ref.MethodByName("TableName")
	if method.IsValid() {
		r := method.Call([]reflect.Value{})
		tableName = r[0].String()
	} else {
		return fmt.Errorf("the current model does not have a table name defined")
	}
	return DB.Table(tableName).Where(where, args...).Updates(updateModel).Error
}

func UpdateByWhere(dataModel interface{}, where string, updates map[string]interface{}) error {
	return DB.Model(dataModel).Where(where).Updates(updates).Error
}

func UpdateByModelWhere(whereModel interface{}, updates map[string]interface{}) error {
	return DB.Model(whereModel).Updates(updates).Error
}

func UpdateByModelWhereModel(whereModel interface{}, updateModel interface{}) error {
	return DB.Model(whereModel).Updates(updateModel).Error
}

func UpdateByArgsWhere(dataModel interface{}, where string, args []interface{}, updates map[string]interface{}) error {
	return DB.Model(dataModel).Where(where, args...).Updates(updates).Error
}

func UpdateByField(dataModel interface{}, where interface{}, column, expr string, updates ...interface{}) error {
	//use example:
	//m := UpdateModel{}
	//err := UpdateByArgsField(&m, "id = ?", 1, "value", "value + ?", 1)
	return DB.Model(dataModel).Where(where).Update(column, gorm.Expr(expr, updates...)).Error
}

func UpdateByArgsField(dataModel interface{}, where string, args []interface{}, column, expr string, updates ...interface{}) error {
	//use example:
	//m := UpdateModel{}
	//err := UpdateByArgsField(&m, "id = ?", []interface{}{1}, "value", "value + ?", 1)
	return DB.Model(dataModel).Where(where, args...).Update(column, gorm.Expr(expr, updates...)).Error
}

func UpdateByModelField(whereModel interface{}, column, expr string, updates ...interface{}) error {
	//use example:
	//m := UpdateModel{}
	//m.Id = 1
	//err := UpdateByModelField(&m, "value", "value + ?", 1)
	return DB.Model(whereModel).Update(column, gorm.Expr(expr, updates...)).Error
}

func DeleteHardById(dataModels interface{}) error {
	return DB.Unscoped().Delete(dataModels).Error
}

func DeleteSoftById(dataModels interface{}) error {
	return DB.Delete(dataModels).Error
}

func DeleteHardByWhere(dataModel interface{}, where string, args []interface{}) error {
	return DB.Unscoped().Where(where, args...).Delete(dataModel).Error
}

func DeleteSoftByWhere(dataModel interface{}, where string, args []interface{}) error {
	return DB.Where(where, args...).Delete(dataModel).Error
}

func RetrieveById(whereModel interface{}) error {
	return DB.First(whereModel).Error
}

func RetrieveByWhere(pageSize, pageNo int, dataModel interface{}, order, where string, args []interface{}) (interface{}, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveByWhere(0, 0, &m, "", "id=?", []interface{}{1})
	//dataList := result.(*[]*dbs.UpdateModel)
	//fmt.Println((*data)[0])
	var count int64 = 0
	if err := DB.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
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
	if err := DB.
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

func RetrieveByWhereString(pageSize, pageNo int, dataModel interface{}, order, where string, args []interface{}) (string, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveByWhereString(0, 0, &m, "", "id=?", []interface{}{1})
	//fmt.Println(result)
	var count int64 = 0
	if err := DB.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
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
	if err := DB.
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

func RetrieveByWhereBytes(pageSize, pageNo int, dataModel interface{}, order, where string, args []interface{}) ([]byte, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveByWhereBytes(0, 0, &m, "", "id=?", []interface{}{1})
	//var dataList []interface{}
	//json.Unmarshal(result,&dataList)
	//fmt.Println(dataList)
	var count int64 = 0
	if err := DB.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
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
	if err := DB.
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

func RetrieveByWhereSelect(pageSize, pageNo int, dataModel interface{}, fields []string, order, where string, args []interface{}) (interface{}, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveBySelect(0, 0, &m, []string{"id", "name"}, "", "id=?", []interface{}{1})
	//dataList := result.(*[]*dbs.UpdateModel)
	//fmt.Println((*dataList)[0])
	var count int64
	if err := DB.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
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
	if err := DB.
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

func RetrieveByWhereSelectString(pageSize, pageNo int, dataModel interface{}, fields []string, order, where string, args []interface{}) (string, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveBySelectString(0, 0, &m, []string{"id", "name"}, "", "id=?", []interface{}{1})
	//fmt.Println(result)
	var count int64
	if err := DB.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
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
	if err := DB.
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

func RetrieveByWhereSelectBytes(pageSize, pageNo int, dataModel interface{}, fields []string, order, where string, args []interface{}) ([]byte, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//result, count, err := RetrieveBySelectBytes(0, 0, &m, []string{"id", "name"}, "", "id=?", []interface{}{1})
	//var dataList []interface{}
	//json.Unmarshal(result,&dataList)
	//fmt.Println(dataList)
	var count int64
	if err := DB.Model(dataModel).Where(where, args...).Count(&count).Error; err != nil {
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
	if err := DB.
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

func RawSql(sql string, values ...interface{}) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)
	rows, err := DB.Raw(sql, values...).Rows()
	if err != nil {
		return results, err
	}
	results = Rows2Map(rows)
	return results, nil
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
				record[colType.Name()] = ""
			} else {
				val, err := json.Marshal(rowValue[i])
				if err != nil {
					record[colType.Name()] = rowValue[i]
				} else {
					record[colType.Name()] = utils.Byte2Any(val, colType.ScanType())
				}
			}
		}
		res = append(res, record)
	}
	return res
}

func Exec(sql string, values ...interface{}) error {
	return DB.Exec(sql, values...).Error
}

func RetrieveByModel(pageSize, pageNo int, whereModel interface{}, order string) (interface{}, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModel(0, 0, &m, "")
	//dataList := result.(*[]*dbs.UpdateModel)
	//fmt.Println((*data)[0])
	var count int64 = 0
	if err := DB.Model(whereModel).Count(&count).Error; err != nil {
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
	if err := DB.
		Model(whereModel).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	return results.Interface(), count, nil
}

func RetrieveByModelString(pageSize, pageNo int, whereModel interface{}, order string) (string, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelString(0, 0, &m, "")
	//fmt.Println(result)
	var count int64 = 0
	if err := DB.Model(whereModel).Count(&count).Error; err != nil {
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
	if err := DB.
		Model(whereModel).
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

func RetrieveByModelBytes(pageSize, pageNo int, whereModel interface{}, order string) ([]byte, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelBytes(0, 0, &m, "")
	//var dataList []interface{}
	//json.Unmarshal(result,&dataList)
	//fmt.Println(dataList)
	var count int64 = 0
	if err := DB.Model(whereModel).Count(&count).Error; err != nil {
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
	if err := DB.
		Model(whereModel).
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

func RetrieveByModelSelect(pageSize, pageNo int, whereModel interface{}, fields []string, order string) (interface{}, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelSelect(0, 0, &m, []string{"id", "name"}, "")
	//dataList := result.(*[]*dbs.UpdateModel)
	//fmt.Println((*dataList)[0])
	var count int64
	if err := DB.Model(whereModel).Count(&count).Error; err != nil {
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
	if err := DB.
		Model(whereModel).
		Select(fields).
		Offset(offset).
		Limit(pageSize).
		Order(order).
		Find(results.Interface()).Error; err != nil {
		return nil, count, err
	}
	return results.Interface(), count, nil
}

func RetrieveByModelSelectString(pageSize, pageNo int, whereModel interface{}, fields []string, order string) (string, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelSelectString(0, 0, &m, []string{"id", "name"}, "")
	//fmt.Println(result)
	var count int64
	if err := DB.Model(whereModel).Count(&count).Error; err != nil {
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
	if err := DB.
		Model(whereModel).
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

func RetrieveByModelSelectBytes(pageSize, pageNo int, whereModel interface{}, fields []string, order string) ([]byte, int64, error) {
	//use example:
	//m := dbs.UpdateModel{}
	//m.Id = 1
	//result, count, err := RetrieveByModelSelectBytes(0, 0, &m, []string{"id", "name"}, "")
	//var dataList []interface{}
	//json.Unmarshal(result,&dataList)
	//fmt.Println(dataList)
	var count int64
	if err := DB.Model(whereModel).Count(&count).Error; err != nil {
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
	if err := DB.
		Model(whereModel).
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

func Transaction(fc func(tx *TX) error) error {
	f := HandleFunc(fc)
	return DB.Transaction(f)
}

func HandleFunc(handler func(tx *TX) error) func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		return handler(&TX{db})
	}
}

func TruncateTable(dataModel interface{}) error {
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
	_, err := RawSql(rSql)
	return err
}
