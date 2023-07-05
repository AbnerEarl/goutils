package dbs

import (
	"database/sql"
	"encoding/json"
	"github.com/YouAreOnlyOne/goutils/utils"
	"gorm.io/gorm"
	"reflect"
	"time"
)

var (
	DefaultLimit    = 10
	DefaultMaxLimit = 1000
)

type BaseModel struct {
	Id        uint64     `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id;comment:'主键ID'"`
	IsDel     uint64     `json:"-" gorm:"column:is_del;default:0;comment:'删除标志'"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;comment:'创建时间'"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at;comment:'更新时间'"`
	DeletedAt *time.Time `json:"-" gorm:"column:deleted_at;comment:'删除时间'" sql:"index"`
	Remark    string     `json:"remark" gorm:"column:remark;null;type:text;comment:'备注信息'"`
}

func (c *BaseModel) TableName() string {
	return "base_info"
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

func UpdateByWhere(dataModel interface{}, where string, updates map[string]interface{}) error {
	return DB.Model(dataModel).Where(where).Updates(updates).Error
}

func UpdateByModelWhere(whereModel interface{}, updates map[string]interface{}, ) error {
	return DB.Model(whereModel).Updates(updates).Error
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

func RetrieveById(id uint64, dataModel interface{}) error {
	return DB.Where("id = ?", id).First(dataModel).Error
}

func RetrieveByModel(whereModel interface{}) error {
	return DB.First(whereModel).Error
}

func RetrieveByWhere(pageSize, pageNo int, dataModel interface{}, where, order string, args []interface{}) (interface{}, int64, error) {
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
	//dataList:=bySelect.(*[]*dbs.UpdateModel)
	//fmt.Println((*data)[0])
	return results.Interface(), count, nil
}

func RetrieveByWhereString(pageSize, pageNo int, dataModel interface{}, where, order string, args []interface{}) (string, int64, error) {
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
	//var dataList []interface{}
	//json.Unmarshal([]byte(string(bytes)),&dataList)
	//fmt.Println(dataList)
	return string(bytes), count, nil
}

func RetrieveByWhereBytes(pageSize, pageNo int, dataModel interface{}, where, order string, args []interface{}) ([]byte, int64, error) {
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
	//var dataList []interface{}
	//json.Unmarshal(bytes,&dataList)
	//fmt.Println(dataList)
	return bytes, count, nil
}

func RetrieveBySelect(pageSize, pageNo int, dataModel interface{}, fields []string, where, order string, args []interface{}) (interface{}, int64, error) {
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
	//dataList:=bySelect.(*[]*dbs.UpdateModel)
	//fmt.Println((*data)[0])
	return results.Interface(), count, nil
}

func RetrieveBySelectString(pageSize, pageNo int, dataModel interface{}, fields []string, where, order string, args []interface{}) (string, int64, error) {
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
	//var dataList []interface{}
	//json.Unmarshal([]byte(string(bytes)),&dataList)
	//fmt.Println(dataList)
	return string(bytes), count, nil
}

func RetrieveBySelectBytes(pageSize, pageNo int, dataModel interface{}, fields []string, where, order string, args []interface{}) ([]byte, int64, error) {
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
	//var dataList []interface{}
	//json.Unmarshal(bytes,&dataList)
	//fmt.Println(dataList)
	return bytes, count, nil
}

func RawSql(sql string) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)
	rows, err := DB.Raw(sql).Rows()
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
				record[colType.Name()] = utils.Byte2Any(rowValue[i].([]byte), colType.ScanType())
			}
		}
		res = append(res, record)
	}
	return res
}
