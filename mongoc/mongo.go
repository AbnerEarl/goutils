package mongoc

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"reflect"
	"time"
)

var (
	MongoCli *mongo.Client
	DbName   string

	CtxExpireTime = 5 * time.Second

	DefaultLimit    int64 = 10
	DefaultMaxLimit int64 = 1000
)

func InitMongo(username, password, host string, port uint64, dbName string, minPoolSize, maxPoolSize uint64) error {
	var mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	ops := options.Client().ApplyURI(mongoURI)
	if minPoolSize < 1 {
		minPoolSize = 100
	}
	if maxPoolSize < 1 {
		maxPoolSize = 10000
	}
	ops.SetMinPoolSize(minPoolSize)
	ops.SetMaxPoolSize(maxPoolSize)

	client, err := mongo.Connect(ctx, ops)
	if err != nil {
		return err
	}
	MongoCli = client
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	//client.Database(dbName)
	DbName = dbName
	return nil
}

func FiledMethodValue(methodName string, dataModel interface{}) (string, error) {
	tableName := ""
	ref := reflect.ValueOf(dataModel)
	switch reflect.TypeOf(dataModel).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(dataModel)
		if s.Len() > 0 {
			ref = s.Index(0)
		}
	case reflect.Ptr:
		s := reflect.ValueOf(dataModel).Elem()
		if s.Len() > 0 {
			ref = s.Index(0)
		}
	}

	method := ref.MethodByName(methodName)
	if method.IsValid() {
		r := method.Call([]reflect.Value{})
		tableName = r[0].String()
		return tableName, nil
	} else {
		return tableName, fmt.Errorf("the current model does not have a collection name defined")
	}
}

func InsertOne(coll string, doc interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := MongoCli.Database(DbName).Collection(coll).InsertOne(ctx, doc)
	return err
}

func InsertMany(coll string, docs []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := MongoCli.Database(DbName).Collection(coll).InsertMany(ctx, docs)
	return err
}

func UpdateOne(coll string, filter, doc interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := MongoCli.Database(DbName).Collection(coll).UpdateOne(ctx, filter, doc)
	return err
}

func UpdateMany(coll string, filter, docs []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := MongoCli.Database(DbName).Collection(coll).UpdateMany(ctx, filter, docs)
	return err
}

func UpdateById(coll string, id, doc interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := MongoCli.Database(DbName).Collection(coll).UpdateByID(ctx, id, doc)
	return err
}

func CreateOne(coll string, indexes []string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	keys := bsonx.Doc{}
	for _, k := range indexes {
		keys = keys.Append(k, bsonx.Int32(1))
	}

	indexName, err := MongoCli.Database(DbName).Collection(coll).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true),
	})
	return indexName, err
}

func CountDocuments(coll string, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	count, err := MongoCli.Database(DbName).Collection(coll).CountDocuments(ctx, filter)
	return count, err
}

func CreateMany(coll string, indexes [][]string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	var indexModels []mongo.IndexModel
	for _, index := range indexes {
		keys := bsonx.Doc{}
		for _, k := range index {
			keys = keys.Append(k, bsonx.Int32(1))
		}
		indexModels = append(indexModels, mongo.IndexModel{
			Keys:    keys,
			Options: options.Index().SetUnique(true),
		})
	}

	indexNames, err := MongoCli.Database(DbName).Collection(coll).Indexes().CreateMany(ctx, indexModels)
	return indexNames, err
}

func FindOne(coll string, filter, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	err := MongoCli.Database(DbName).Collection(coll).FindOne(ctx, filter).Decode(result)
	return err
}

func FindMany(pageSize, pageNo int64, coll string, filter, sort, result, projects interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	ops := options.Find()
	ops.SetSort(sort)
	ops.SetLimit(pageSize)
	ops.SetSkip(offset)
	ops.SetProjection(projects)
	cur, err := MongoCli.Database(DbName).Collection(coll).Find(ctx, filter, ops)
	if err != nil {
		return err
	}
	err = cur.All(ctx, result)
	return err
}

func FindManyCount(pageSize, pageNo int64, coll string, filter, sort, result, projects interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	ops := options.Find()
	ops.SetSort(sort)
	ops.SetLimit(pageSize)
	ops.SetSkip(offset)
	ops.SetProjection(projects)
	cur, err := MongoCli.Database(DbName).Collection(coll).Find(ctx, filter, ops)
	if err != nil {
		return 0, err
	}
	err = cur.All(ctx, result)
	if err != nil {
		return 0, err
	}
	count, err := CountDocuments(coll, filter)
	return count, err
}

func FindManyRange(pageSize, lastId int64, coll string, filter, sort, result, projects interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	ops := options.Find()
	ops.SetSort(sort)
	ops.SetLimit(pageSize)
	ops.SetProjection(projects)
	page := bson.M{"_id": bson.M{"$gt": lastId}}
	cur, err := MongoCli.Database(DbName).Collection(coll).Find(ctx, bson.A{page, filter}, ops)
	if err != nil {
		return err
	}
	err = cur.All(ctx, result)
	return err
}

func FindManyRangeCount(pageSize, lastId int64, coll string, filter, sort, result, projects interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	ops := options.Find()
	ops.SetSort(sort)
	ops.SetLimit(pageSize)
	ops.SetProjection(projects)
	page := bson.M{"_id": bson.M{"$gt": lastId}}
	cur, err := MongoCli.Database(DbName).Collection(coll).Find(ctx, bson.A{page, filter}, ops)
	if err != nil {
		return 0, err
	}
	err = cur.All(ctx, result)
	if err != nil {
		return 0, err
	}
	count, err := CountDocuments(coll, filter)
	return count, err
}

func DeleteOne(coll string, filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := MongoCli.Database(DbName).Collection(coll).DeleteOne(ctx, filter)
	return err
}

func DeleteMany(coll string, filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err := MongoCli.Database(DbName).Collection(coll).DeleteMany(ctx, filter)
	return err
}

func FindFindOneAndDeleteOne(coll string, filter, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	err := MongoCli.Database(DbName).Collection(coll).FindOneAndDelete(ctx, filter).Decode(result)
	return err
}

func FindOneAndReplace(coll string, filter, replace, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	err := MongoCli.Database(DbName).Collection(coll).FindOneAndReplace(ctx, filter, replace).Decode(result)
	return err
}

func FindOneAndUpdate(coll string, filter, update, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	err := MongoCli.Database(DbName).Collection(coll).FindOneAndUpdate(ctx, filter, update).Decode(result)
	return err
}

func Transaction(fn func(mongo.SessionContext) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	err := MongoCli.UseSession(ctx, fn)
	return err
}

func TransactionWithOptions(opts *options.SessionOptions, fn func(mongo.SessionContext) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	err := MongoCli.UseSessionWithOptions(ctx, opts, fn)
	return err
}

func Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	err := MongoCli.Disconnect(ctx)
	return err
}

func InsertOneByModel(doc interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	coll, err := FiledMethodValue("CollectionName", doc)
	if err != nil {
		return err
	}
	_, err = MongoCli.Database(DbName).Collection(coll).InsertOne(ctx, doc)
	return err
}

func InsertManyByModel(docs []interface{}) error {
	coll, err := FiledMethodValue("CollectionName", docs[0])
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err = MongoCli.Database(DbName).Collection(coll).InsertMany(ctx, docs)
	return err
}

func UpdateOneByModel(filter, doc interface{}) error {
	coll, err := FiledMethodValue("CollectionName", doc)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err = MongoCli.Database(DbName).Collection(coll).UpdateOne(ctx, filter, doc)
	return err
}

func UpdateManyByModel(filter, docs []interface{}) error {
	coll, err := FiledMethodValue("CollectionName", docs[0])
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err = MongoCli.Database(DbName).Collection(coll).UpdateMany(ctx, filter, docs)
	return err
}

func UpdateByIdByModel(id, doc interface{}) error {
	coll, err := FiledMethodValue("CollectionName", doc)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err = MongoCli.Database(DbName).Collection(coll).UpdateByID(ctx, id, doc)
	return err
}

func CreateOneByModel(doc interface{}, indexes []string) (string, error) {
	coll, err := FiledMethodValue("CollectionName", doc)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	keys := bsonx.Doc{}
	for _, k := range indexes {
		keys = keys.Append(k, bsonx.Int32(1))
	}

	indexName, err := MongoCli.Database(DbName).Collection(coll).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true),
	})
	return indexName, err
}

func CountDocumentsByModel(doc, filter interface{}) (int64, error) {
	coll, err := FiledMethodValue("CollectionName", doc)
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	count, err := MongoCli.Database(DbName).Collection(coll).CountDocuments(ctx, filter)
	return count, err
}

func CreateManyByModel(doc interface{}, indexes [][]string) ([]string, error) {
	coll, err := FiledMethodValue("CollectionName", doc)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	var indexModels []mongo.IndexModel
	for _, index := range indexes {
		keys := bsonx.Doc{}
		for _, k := range index {
			keys = keys.Append(k, bsonx.Int32(1))
		}
		indexModels = append(indexModels, mongo.IndexModel{
			Keys:    keys,
			Options: options.Index().SetUnique(true),
		})
	}

	indexNames, err := MongoCli.Database(DbName).Collection(coll).Indexes().CreateMany(ctx, indexModels)
	return indexNames, err
}

func FindOneByModel(filter, result interface{}) error {
	coll, err := FiledMethodValue("CollectionName", result)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	err = MongoCli.Database(DbName).Collection(coll).FindOne(ctx, filter).Decode(result)
	return err
}

func FindManyByModel(pageSize, pageNo int64, coll string, filter, sort, result, projects interface{}) error {
	coll, err := FiledMethodValue("CollectionName", result)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	ops := options.Find()
	ops.SetSort(sort)
	ops.SetLimit(pageSize)
	ops.SetSkip(offset)
	ops.SetProjection(projects)
	cur, err := MongoCli.Database(DbName).Collection(coll).Find(ctx, filter, ops)
	if err != nil {
		return err
	}
	err = cur.All(ctx, &result)
	return err
}

func FindManyCountByModel(pageSize, pageNo int64, coll string, filter, sort, result, projects interface{}) (int64, error) {
	coll, err := FiledMethodValue("CollectionName", result)
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	if pageNo == 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageSize
	ops := options.Find()
	ops.SetSort(sort)
	ops.SetLimit(pageSize)
	ops.SetSkip(offset)
	ops.SetProjection(projects)
	cur, err := MongoCli.Database(DbName).Collection(coll).Find(ctx, filter, ops)
	if err != nil {
		return 0, err
	}
	err = cur.All(ctx, result)
	if err != nil {
		return 0, err
	}
	count, err := CountDocuments(coll, filter)
	return count, err
}

func FindManyRangeByModel(pageSize, lastId int64, coll string, filter, sort, result, projects interface{}) error {
	coll, err := FiledMethodValue("CollectionName", result)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	ops := options.Find()
	ops.SetSort(sort)
	ops.SetLimit(pageSize)
	ops.SetProjection(projects)
	page := bson.M{"_id": bson.M{"$gt": lastId}}
	cur, err := MongoCli.Database(DbName).Collection(coll).Find(ctx, bson.A{page, filter}, ops)
	if err != nil {
		return err
	}
	err = cur.All(ctx, result)
	return err
}

func FindManyRangeCountByModel(pageSize, lastId int64, coll string, filter, sort, result, projects interface{}) (int64, error) {
	coll, err := FiledMethodValue("CollectionName", result)
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	if pageSize == 0 {
		pageSize = DefaultLimit
	} else if pageSize > DefaultMaxLimit {
		pageSize = DefaultMaxLimit
	}
	ops := options.Find()
	ops.SetSort(sort)
	ops.SetLimit(pageSize)
	ops.SetProjection(projects)
	page := bson.M{"_id": bson.M{"$gt": lastId}}
	cur, err := MongoCli.Database(DbName).Collection(coll).Find(ctx, bson.A{page, filter}, ops)
	if err != nil {
		return 0, err
	}
	err = cur.All(ctx, result)
	if err != nil {
		return 0, err
	}
	count, err := CountDocuments(coll, filter)
	return count, err
}

func DeleteOneByModel(doc, filter interface{}) error {
	coll, err := FiledMethodValue("CollectionName", doc)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err = MongoCli.Database(DbName).Collection(coll).DeleteOne(ctx, filter)
	return err
}

func DeleteManyByModel(doc, filter interface{}) error {
	coll, err := FiledMethodValue("CollectionName", doc)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	_, err = MongoCli.Database(DbName).Collection(coll).DeleteMany(ctx, filter)
	return err
}

func FindFindOneAndDeleteOneByModel(filter, result interface{}) error {
	coll, err := FiledMethodValue("CollectionName", result)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	err = MongoCli.Database(DbName).Collection(coll).FindOneAndDelete(ctx, filter).Decode(result)
	return err
}

func FindOneAndReplaceByModel(filter, replace, result interface{}) error {
	coll, err := FiledMethodValue("CollectionName", result)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	err = MongoCli.Database(DbName).Collection(coll).FindOneAndReplace(ctx, filter, replace).Decode(result)
	return err
}

func FindOneAndUpdateByModel(filter, update, result interface{}) error {
	coll, err := FiledMethodValue("CollectionName", result)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	err = MongoCli.Database(DbName).Collection(coll).FindOneAndUpdate(ctx, filter, update).Decode(result)
	return err
}
