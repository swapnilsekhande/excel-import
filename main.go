package main

import (
	"excel-import/database"
	"excel-import/migrations"
	"excel-import/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Error(err)
	}
}

func main() {
	database.InitBookStore()
	database.InitRedisConnections()
	defer func() {
		if db, err := database.EmployeeDetails.DB(); err == nil {
			logrus.Error(err)
			db.Close()
		}
	}()

	if errors := migrations.MigrationRun(database.EmployeeDetails); errors != nil {
		logrus.Error("Error : ", errors)
	}
	router := gin.Default()
	routes.RenderRoutes(router)
}
