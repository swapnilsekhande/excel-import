package controllers

import (
	"excel-import/database"
	"excel-import/models"
	"excel-import/repositories"
	"excel-import/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/*
* @Description: This function is used to upload the excel file and save it to the server primarily then get Content from the file validate it and save it to the database
* @param c: The gin context
* @return: nil
 */
func UploadXlsxFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	c.SaveUploadedFile(file, "./Uploads/"+file.Filename)
	employeeDetails, err := services.ReadExcelFile("./Uploads/" + file.Filename)
	if err != nil {
		logrus.Error("Error reading Excel file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read Excel file"})
		return
	}

	for _, employee := range employeeDetails {
		if err := database.EmployeeDetails.Create(&employee).Error; err != nil {
			logrus.Error("Error saving employee details to database:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save employee details" + err.Error()})
			return
		}
		repositories.CacheEmployee(employee)
		if err != nil {
			logrus.Error("Error caching employee details:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache employee details"})
			return
		}
		logrus.Info("Employee details saved and cached successfully:", employee)
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

func GetEmployee(c *gin.Context) {
	// Convert ID from string to uint
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	// Fetch employee (from Redis or MySQL)
	employee, err := services.GetEmployeeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func UpdateEmployee(c *gin.Context) {
	var employee models.EmployeeDetails

	// Bind JSON from the request body to the employee struct
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Call the update function
	if err := repositories.UpdateEmployee(employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully"})
}

func GetPaginatedEmployees(c *gin.Context) {
	// Get query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err1 := strconv.Atoi(pageStr)
	limit, err2 := strconv.Atoi(limitStr)

	if err1 != nil || err2 != nil || page < 1 || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	employees, err := repositories.GetAllEmployees(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employees", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"limit":     limit,
		"employees": employees,
	})
}
