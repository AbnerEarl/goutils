package dbs

import (
	"fmt"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

var (
	DefaultMaxConn  uint64 = 2000
	DefaultIdleConn uint64 = 20
)

func OpenDBMySQL(username, password, ip string, port uint64, dbName string, dryRun bool, maxConn, idleConn uint64, customConn string) *DB {
	if len(ip) < 4 {
		ip = "127.0.0.1"
	}
	if port < 1 {
		port = 3306
	}
	if maxConn < 1 {
		maxConn = DefaultMaxConn
	}
	if idleConn < 1 {
		idleConn = DefaultIdleConn
	}

	dsn := customConn
	if len(dsn) < 4 {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, ip, port, dbName)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DryRun: dryRun,
	})
	if err != nil {
		panic(err)
	}

	dc, _ := db.DB()
	dc.SetMaxOpenConns(int(maxConn))  // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	dc.SetMaxIdleConns(int(idleConn)) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	return &DB{db}
}

func OpenDBPostgreSQL(username, password, ip string, port uint64, dbName string, dryRun bool, maxConn, idleConn uint64, customConn string) *DB {
	if len(ip) < 4 {
		ip = "127.0.0.1"
	}
	if port < 1 {
		port = 9920
	}
	if maxConn < 1 {
		maxConn = DefaultMaxConn
	}
	if idleConn < 1 {
		idleConn = DefaultIdleConn
	}

	dsn := customConn
	if len(dsn) < 4 {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", ip, username, password, dbName, port)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DryRun: dryRun,
	})
	if err != nil {
		panic(err)
	}

	dc, _ := db.DB()
	dc.SetMaxOpenConns(int(maxConn))  // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	dc.SetMaxIdleConns(int(idleConn)) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	return &DB{db}
}

func OpenDBSQLite(fileName string, dryRun bool, maxConn, idleConn uint64, customConn string) *DB {
	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{
		DryRun: dryRun,
	})
	if err != nil {
		panic(err)
	}

	dc, _ := db.DB()
	dc.SetMaxOpenConns(int(maxConn))  // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	dc.SetMaxIdleConns(int(idleConn)) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	return &DB{db}
}

func OpenDBSQLServer(username, password, ip string, port uint64, dbName string, dryRun bool, maxConn, idleConn uint64, customConn string) *DB {
	if len(ip) < 4 {
		ip = "127.0.0.1"
	}
	if port < 1 {
		port = 9930
	}
	if maxConn < 1 {
		maxConn = DefaultMaxConn
	}
	if idleConn < 1 {
		idleConn = DefaultIdleConn
	}

	dsn := customConn
	if len(dsn) < 4 {
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", username, password, ip, port, dbName)
	}
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		DryRun: dryRun,
	})
	if err != nil {
		panic(err)
	}

	dc, _ := db.DB()
	dc.SetMaxOpenConns(int(maxConn))  // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	dc.SetMaxIdleConns(int(idleConn)) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	return &DB{db}
}

func OpenDBTiDB(username, password, ip string, port uint64, dbName string, dryRun bool, maxConn, idleConn uint64, customConn string) *DB {
	if len(ip) < 4 {
		ip = "127.0.0.1"
	}
	if port < 1 {
		port = 4000
	}
	if maxConn < 1 {
		maxConn = DefaultMaxConn
	}
	if idleConn < 1 {
		idleConn = DefaultIdleConn
	}

	dsn := customConn
	if len(dsn) < 4 {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&tls=register-tidb-tls", username, password, ip, port, dbName)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DryRun: dryRun,
	})
	if err != nil {
		panic(err)
	}

	dc, _ := db.DB()
	dc.SetMaxOpenConns(int(maxConn))  // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	dc.SetMaxIdleConns(int(idleConn)) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	return &DB{db}
}

func OpenDBClickhouse(username, password, ip string, port uint64, dbName string, dryRun bool, maxConn, idleConn uint64, customConn string) *DB {
	if len(ip) < 4 {
		ip = "127.0.0.1"
	}
	if port < 1 {
		port = 9000
	}
	if maxConn < 1 {
		maxConn = DefaultMaxConn
	}
	if idleConn < 1 {
		idleConn = DefaultIdleConn
	}

	dsn := customConn
	if len(dsn) < 4 {
		dsn = fmt.Sprintf("tcp://%s:%d?database=%s&username=%s&password=%s&read_timeout=10&write_timeout=20", ip, port, dbName, username, password)
	}
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{
		DryRun: dryRun,
	})
	if err != nil {
		panic(err)
	}

	dc, _ := db.DB()
	dc.SetMaxOpenConns(int(maxConn))  // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	dc.SetMaxIdleConns(int(idleConn)) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	return &DB{db}
}

func (db *DB) Migration(models []interface{}) {
	for _, md := range models {
		db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(md)
	}
}
