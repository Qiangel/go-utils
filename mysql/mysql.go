package mysql

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Username       string `yaml:"user"`
	Password       string
	Host           string `yaml:"host" default:"127.0.0.1"`
	Port           string `default:"3306"`
	Dbname         string `yaml:"db_name"`
	ConnectTimeout int    `default:"10"`
	Debug          bool   `default:"true"`
	MaxLifetime    int    `default:"7200"` //设置连接可以重用的最长时间(单位：秒)
	MaxOpenConns   int    `default:"150"`  //设置数据库的最大打开连接数
	MaxIdleConns   int    `default:"50"`   //设置空闲连接池中的最大连接数
}

func Open(cfg Config) (*gorm.DB, error) {
	dsn := cfg.Username + ":" + cfg.Password +
		"@tcp(" + cfg.Host + ":" + cfg.Port + ")/" + cfg.Dbname + "?" +
		"charset=utf8mb4&parseTime=True&loc=Local"
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig)
	if err != nil {
		return nil, err
	}
	if cfg.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)
	return db, nil
}
