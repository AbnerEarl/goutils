package mongoc

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
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

func FinMany(pageSize, pageNo int64, coll string, filter, sort, result, projects interface{}) error {
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

func FinManyCount(pageSize, pageNo int64, coll string, filter, sort, result, projects interface{}) (int64, error) {
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

func Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), CtxExpireTime)
	defer cancel()
	err := MongoCli.Disconnect(ctx)
	return err
}

type BaseModel struct {
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"-" bson:"deleted_at"`
	Remark    string     `json:"remark" bson:"remark"`
}
