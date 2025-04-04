package routes

import (
	"excel-import/controllers"

	"github.com/gin-gonic/gin"
)

func RenderRoutes(router *gin.Engine) {
	routerBooksGroup := router.Group("/employee")
	{
		routerBooksGroup.POST("/upload", controllers.UploadXlsxFile)
		routerBooksGroup.GET("/get/:id", controllers.GetEmployee)
		routerBooksGroup.PUT("/update", controllers.UpdateEmployee)
		routerBooksGroup.GET("/list", controllers.GetEmployeesList)
	}
	router.Run(":9000")
}
