package dbs

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

func OpenDBMySQLCluster(dsns []string, dryRun bool, maxConn, idleConn uint64) *DB {
	/**
	@param dsns, such as: []string{Db1Dsn, Db2Dsn, Db3Dsn, Db4Dsn}
	var Db1Dsn = "root:password@tcp(localhost:3306)/db_ex1?charset=utf8&parseTime=true"
	var Db2Dsn = "root:password@tcp(localhost:3306)/db_ex2?charset=utf8&parseTime=true"
	var Db3Dsn = "root:password@tcp(localhost:3306)/db_ex3?charset=utf8&parseTime=true"
	var Db4Dsn = "root:password@tcp(localhost:3306)/db_ex4?charset=utf8&parseTime=true"
	*/

	number := len(dsns)
	if number < 1 {
		return nil
	}

	if maxConn < 1 {
		maxConn = DefaultMaxConn
	}
	if idleConn < 1 {
		idleConn = DefaultIdleConn
	}

	db, err := gorm.Open(mysql.Open(dsns[0]), &gorm.Config{
		DryRun: dryRun,
	})
	if err != nil {
		panic(err)
	}

	if number > 1 {
		if number < 3 {
			var sources []gorm.Dialector
			for i := 1; i < number; i++ {
				sources = append(sources, mysql.Open(dsns[i]))
			}
			db.Use(dbresolver.Register(dbresolver.Config{
				Sources: sources,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}))
		} else {
			n := number/2 + 1
			var sources []gorm.Dialector
			for i := 1; i < n; i++ {
				sources = append(sources, mysql.Open(dsns[i]))
			}
			var replicas []gorm.Dialector
			for i := n; i < number; i++ {
				replicas = append(replicas, mysql.Open(dsns[i]))
			}
			db.Use(dbresolver.Register(dbresolver.Config{
				Sources:  sources,
				Replicas: replicas,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}))
		}
	}

	dc, _ := db.DB()
	dc.SetMaxOpenConns(int(maxConn))  // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接出现too many connections的错误。
	dc.SetMaxIdleConns(int(idleConn)) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	dc.SetConnMaxIdleTime(time.Hour)
	dc.SetConnMaxLifetime(24 * time.Hour)
	return &DB{db}
}

func OpenDBPostgreSQLCluster(dsns []string, dryRun bool, maxConn, idleConn uint64) *DB {
	/**
	@param dsns, such as: []string{Db1Dsn, Db2Dsn, Db3Dsn, Db4Dsn}
	var Db1Dsn = "host=localhost user=root password=root dbname=db_ex1 port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	var Db2Dsn = "host=localhost user=root password=root dbname=db_ex2 port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	var Db3Dsn = "host=localhost user=root password=root dbname=db_ex3 port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	var Db4Dsn = "host=localhost user=root password=root dbname=db_ex4 port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	*/

	number := len(dsns)
	if number < 1 {
		return nil
	}
	if maxConn < 1 {
		maxConn = DefaultMaxConn
	}
	if idleConn < 1 {
		idleConn = DefaultIdleConn
	}

	db, err := gorm.Open(postgres.Open(dsns[0]), &gorm.Config{
		DryRun: dryRun,
	})
	if err != nil {
		panic(err)
	}

	if number > 1 {
		if number < 3 {
			var sources []gorm.Dialector
			for i := 1; i < number; i++ {
				sources = append(sources, postgres.Open(dsns[i]))
			}
			db.Use(dbresolver.Register(dbresolver.Config{
				Sources: sources,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}))
		} else {
			n := number/2 + 1
			var sources []gorm.Dialector
			for i := 1; i < n; i++ {
				sources = append(sources, postgres.Open(dsns[i]))
			}
			var replicas []gorm.Dialector
			for i := n; i < number; i++ {
				replicas = append(replicas, postgres.Open(dsns[i]))
			}
			db.Use(dbresolver.Register(dbresolver.Config{
				Sources:  sources,
				Replicas: replicas,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}))
		}
	}

	dc, _ := db.DB()
	dc.SetMaxOpenConns(int(maxConn))  // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接出现too many connections的错误。
	dc.SetMaxIdleConns(int(idleConn)) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	dc.SetConnMaxIdleTime(time.Hour)
	dc.SetConnMaxLifetime(24 * time.Hour)
	return &DB{db}
}

func OpenDBSQLiteCluster(fileNames []string, dryRun bool, maxConn, idleConn uint64, customConn string) *DB {
	/**
	@param fileNames, such as: []string{fileName1, fileName2, fileName3, fileName4}
	var fileName1 = "/root/db_ex1.db"
	var fileName2 = "/root/db_ex2.db"
	var fileName3 = "/root/db_ex3.db"
	var fileName4 = "/root/db_ex4.db"
	*/

	number := len(fileNames)
	if number < 1 {
		return nil
	}
	db, err := gorm.Open(sqlite.Open(fileNames[0]), &gorm.Config{
		DryRun: dryRun,
	})
	if err != nil {
		panic(err)
	}

	if number > 1 {
		if number < 3 {
			var sources []gorm.Dialector
			for i := 1; i < number; i++ {
				sources = append(sources, sqlite.Open(fileNames[i]))
			}
			db.Use(dbresolver.Register(dbresolver.Config{
				Sources: sources,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}))
		} else {
			n := number/2 + 1
			var sources []gorm.Dialector
			for i := 1; i < n; i++ {
				sources = append(sources, sqlite.Open(fileNames[i]))
			}
			var replicas []gorm.Dialector
			for i := n; i < number; i++ {
				replicas = append(replicas, sqlite.Open(fileNames[i]))
			}
			db.Use(dbresolver.Register(dbresolver.Config{
				Sources:  sources,
				Replicas: replicas,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}))
		}
	}

	dc, _ := db.DB()
	dc.SetMaxOpenConns(int(maxConn))  // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接出现too many connections的错误。
	dc.SetMaxIdleConns(int(idleConn)) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	dc.SetConnMaxIdleTime(time.Hour)
	dc.SetConnMaxLifetime(24 * time.Hour)
	return &DB{db}
}

func OpenDBSQLServerCluster(dsns []string, dryRun bool, maxConn, idleConn uint64) *DB {
	/**
	@param dsns, such as: []string{Db1Dsn, Db2Dsn, Db3Dsn, Db4Dsn}
	var Db1Dsn = "sqlserver://root:password@localhost:9930?database=db_ex1"
	var Db2Dsn = "sqlserver://root:password@localhost:9930?database=db_ex2"
	var Db3Dsn = "sqlserver://root:password@localhost:9930?database=db_ex3"
	var Db4Dsn = "sqlserver://root:password@localhost:9930?database=db_ex4"
	*/

	number := len(dsns)
	if number < 1 {
		return nil
	}
	if maxConn < 1 {
		maxConn = DefaultMaxConn
	}
	if idleConn < 1 {
		idleConn = DefaultIdleConn
	}

	db, err := gorm.Open(sqlserver.Open(dsns[0]), &gorm.Config{
		DryRun: dryRun,
	})
	if err != nil {
		panic(err)
	}

	if number > 1 {
		if number < 3 {
			var sources []gorm.Dialector
			for i := 1; i < number; i++ {
				sources = append(sources, sqlserver.Open(dsns[i]))
			}
			db.Use(dbresolver.Register(dbresolver.Config{
				Sources: sources,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}))
		} else {
			n := number/2 + 1
			var sources []gorm.Dialector
			for i := 1; i < n; i++ {
				sources = append(sources, sqlserver.Open(dsns[i]))
			}
			var replicas []gorm.Dialector
			for i := n; i < number; i++ {
				replicas = append(replicas, sqlserver.Open(dsns[i]))
			}
			db.Use(dbresolver.Register(dbresolver.Config{
				Sources:  sources,
				Replicas: replicas,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}))
		}
	}

	dc, _ := db.DB()
	dc.SetMaxOpenConns(int(maxConn))  // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接出现too many connections的错误。
	dc.SetMaxIdleConns(int(idleConn)) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	dc.SetConnMaxIdleTime(time.Hour)
	dc.SetConnMaxLifetime(24 * time.Hour)
	return &DB{db}
}

func OpenDBTiDBCluster(dsns []string, dryRun bool, maxConn, idleConn uint64) *DB {
	/**
	@param dsns, such as: []string{Db1Dsn, Db2Dsn, Db3Dsn, Db4Dsn}
	var Db1Dsn = "root:password@tcp(localhost:4000)/db_ex1?charset=utf8&tls=register-tidb-tls"
	var Db2Dsn = "root:password@tcp(localhost:4000)/db_ex2?charset=utf8&tls=register-tidb-tls"
	var Db3Dsn = "root:password@tcp(localhost:4000)/db_ex3?charset=utf8&tls=register-tidb-tls"
	var Db4Dsn = "root:password@tcp(localhost:4000)/db_ex4?charset=utf8&tls=register-tidb-tls"
	*/

	number := len(dsns)
	if number < 1 {
		return nil
	}
	if maxConn < 1 {
		maxConn = DefaultMaxConn
	}
	if idleConn < 1 {
		idleConn = DefaultIdleConn
	}

	db, err := gorm.Open(mysql.Open(dsns[0]), &gorm.Config{
		DryRun: dryRun,
	})
	if err != nil {
		panic(err)
	}

	if number > 1 {
		if number < 3 {
			var sources []gorm.Dialector
			for i := 1; i < number; i++ {
				sources = append(sources, mysql.Open(dsns[i]))
			}
			db.Use(dbresolver.Register(dbresolver.Config{
				Sources: sources,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}))
		} else {
			n := number/2 + 1
			var sources []gorm.Dialector
			for i := 1; i < n; i++ {
				sources = append(sources, mysql.Open(dsns[i]))
			}
			var replicas []gorm.Dialector
			for i := n; i < number; i++ {
				replicas = append(replicas, mysql.Open(dsns[i]))
			}
			db.Use(dbresolver.Register(dbresolver.Config{
				Sources:  sources,
				Replicas: replicas,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}))
		}
	}

	dc, _ := db.DB()
	dc.SetMaxOpenConns(int(maxConn))  // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接出现too many connections的错误。
	dc.SetMaxIdleConns(int(idleConn)) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	dc.SetConnMaxIdleTime(time.Hour)
	dc.SetConnMaxLifetime(24 * time.Hour)
	return &DB{db}
}

func OpenDBClickhouseCluster(dsns []string, dryRun bool, maxConn, idleConn uint64) *DB {
	/**
	@param dsns, such as: []string{Db1Dsn, Db2Dsn, Db3Dsn, Db4Dsn}
	var Db1Dsn = "tcp://localhost:9000?database=db_ex1&username=root&password=root&read_timeout=10&write_timeout=20"
	var Db2Dsn = "tcp://localhost:9000?database=db_ex2&username=root&password=root&read_timeout=10&write_timeout=20"
	var Db3Dsn = "tcp://localhost:9000?database=db_ex3&username=root&password=root&read_timeout=10&write_timeout=20"
	var Db4Dsn = "tcp://localhost:9000?database=db_ex4&username=root&password=root&read_timeout=10&write_timeout=20"
	*/

	number := len(dsns)
	if number < 1 {
		return nil
	}
	if maxConn < 1 {
		maxConn = DefaultMaxConn
	}
	if idleConn < 1 {
		idleConn = DefaultIdleConn
	}

	db, err := gorm.Open(clickhouse.Open(dsns[0]), &gorm.Config{
		DryRun: dryRun,
	})
	if err != nil {
		panic(err)
	}

	if number > 1 {
		if number < 3 {
			var sources []gorm.Dialector
			for i := 1; i < number; i++ {
				sources = append(sources, clickhouse.Open(dsns[i]))
			}
			db.Use(dbresolver.Register(dbresolver.Config{
				Sources: sources,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}))
		} else {
			n := number/2 + 1
			var sources []gorm.Dialector
			for i := 1; i < n; i++ {
				sources = append(sources, clickhouse.Open(dsns[i]))
			}
			var replicas []gorm.Dialector
			for i := n; i < number; i++ {
				replicas = append(replicas, clickhouse.Open(dsns[i]))
			}
			db.Use(dbresolver.Register(dbresolver.Config{
				Sources:  sources,
				Replicas: replicas,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}))
		}
	}

	dc, _ := db.DB()
	dc.SetMaxOpenConns(int(maxConn))  // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接出现too many connections的错误。
	dc.SetMaxIdleConns(int(idleConn)) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	dc.SetConnMaxIdleTime(time.Hour)
	dc.SetConnMaxLifetime(24 * time.Hour)
	return &DB{db}
}
