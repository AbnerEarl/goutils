package dbs

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func OpenDB(username, password, ip string, port uint64, dbName, dbType string, dryRun bool, maxConn, idleConn uint64) *gorm.DB {
	if len(ip) < 4 {
		ip = "127.0.0.1"
	}
	if port < 1 {
		port = 3306
	}
	if len(dbType) < 2 {
		dbType = "mysql" //mysql, postgres
	}
	if maxConn < 1 {
		maxConn = 20000
	}
	if dbType == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, ip, port, dbName)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			DryRun: dryRun,
		})
		if err != nil {
			panic(err)
		}
		DB = db
	} else if dbType == "postgres" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", ip, username, password, dbName, port)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			DryRun: dryRun,
		})
		if err != nil {
			panic(err)
		}
		DB = db
	} else {
		panic("the database type connection is currently not supported")
	}
	dc, _ := DB.DB()
	dc.SetMaxOpenConns(int(maxConn))  // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	dc.SetMaxIdleConns(int(idleConn)) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	return DB
}

func Migration(models []interface{}) {
	for _, md := range models {
		DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(md)
	}
}

func GetDBClient() *gorm.DB {
	if DB == nil {
		panic("The database has not yet initialized the connection.")
	}
	return DB
}
