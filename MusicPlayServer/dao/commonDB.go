package dao

import (
	"MusicPlayServer/common/config"
	"MusicPlayServer/common/log"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var DBClient *gorm.DB

type Model struct {
	Id        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
}

func InitDB(cfg *config.Config) (*gorm.DB, error) {

	datasource := cfg.Mysql.DataSource
	if len(datasource) == 0 {
		log.Error("Mysql datasource is empty.")
		return nil, errors.New("datasource is empty")
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       datasource, // DSN data source name
		DefaultStringSize:         256,        // string 类型字段的默认长度
		DisableDatetimePrecision:  true,       // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,       // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,       // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,      // 根据当前 MySQL 版本自动配置

	}), &gorm.Config{})

	//db.AutoMigrate(&TUser{}, &Post{})
	db.AutoMigrate(&LikeCountModel{})
	if err != nil {
		panic("failed to connect database")
	}
	// ----------------------------数据库连接池----------------------------
	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	fmt.Println("success to link mysql")
	return db, nil
}
