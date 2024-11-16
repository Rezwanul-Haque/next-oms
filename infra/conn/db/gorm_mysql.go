package db

import (
	"fmt"
	"next-oms/infra/config"
	"next-oms/infra/conn/db/models"
	"next-oms/infra/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type DatabaseClient struct {
	lc logger.LogClient
	DB *gorm.DB
}

func connectMySQL(lc logger.LogClient) {
	conf := config.Db().MySQL

	logger.Client().Info("connecting to mysql at " + conf.Host + ":" + conf.Port + "...")

	logMode := gormlogger.Silent
	if conf.Debug {
		logMode = gormlogger.Info
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", conf.User, conf.Pass, conf.Host, conf.Port, conf.Schema)

	dB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      gormlogger.Default.LogMode(logMode),
	})

	if err != nil {
		panic(err)
	}

	sqlDb, err := dB.DB()
	if err != nil {
		panic(err)
	}

	if conf.MaxIdleConn != 0 {
		sqlDb.SetMaxIdleConns(conf.MaxIdleConn)
	}
	if conf.MaxOpenConn != 0 {
		sqlDb.SetMaxOpenConns(conf.MaxOpenConn)
	}
	if conf.MaxConnLifetime != 0 {
		sqlDb.SetConnMaxLifetime(conf.MaxConnLifetime * time.Second)
	}

	client.DB = dB
	client.lc = lc

	client.DB.AutoMigrate(
		&models.Order{},
		&models.User{},
	)

	logger.Client().Info("mysql connection successful...")
}
