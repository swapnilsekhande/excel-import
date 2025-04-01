package database

import (
	"excel-import/config"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var EmployeeDetails *gorm.DB

func InitBookStore() error {
	dbConfig := config.GetMysqlConfigurationFromEnv()
	var err error
	dsn := dbConfig.DB_UserName + ":" + dbConfig.DB_PassWord + "@tcp(127.0.0.1:" + dbConfig.DB_Port + ")/" + dbConfig.DB_Name + "?charset=utf8mb4&parseTime=True&loc=Local"
	EmployeeDetails, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Error("Failed to connect to the database:", err)
		return err
	}

	// Connection Pool
	// GORM using database/sql to maintain connection pool
	sqlDB, err := EmployeeDetails.DB()
	if err != nil {
		logrus.Error("failed to pool connection: ", err)
		return err
	}
	logrus.Info("Connection pool settings are being configured.")
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}
