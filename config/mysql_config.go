package config

import "os"

type MysqlConfigurationGetEnv struct {
	DB_HostName string
	DB_UserName string
	DB_PassWord string
	DB_Port     string
	DB_Name     string
}

func GetMysqlConfigurationFromEnv() MysqlConfigurationGetEnv {
	return MysqlConfigurationGetEnv{
		DB_HostName: os.Getenv("db_host_name"),
		DB_UserName: os.Getenv("db_user_name"),
		DB_PassWord: os.Getenv("db_password"),
		DB_Port:     os.Getenv("db_port"),
		DB_Name:     os.Getenv("db_database_name"),
	}
}
