package repositories

import (
	"context"
	"encoding/json"
	"excel-import/database"
	"excel-import/models"
	"fmt"
	"time"
)

// Cache key format
func getEmployeeCacheKey(id uint) string {
	return fmt.Sprintf("employee:%d", id)
}

// Save employee data in Redis (cache for 10 minutes)
func CacheEmployee(employee models.EmployeeDetails) error {
	ctx := context.Background()
	data, err := json.Marshal(employee)
	if err != nil {
		return err
	}
	return database.EmployeeDetailsRedisClient.Set(ctx, getEmployeeCacheKey(employee.ID), data, 10*time.Minute).Err()
}

// Retrieve employee data from Redis
func GetCachedEmployee(id uint) (*models.EmployeeDetails, error) {
	ctx := context.Background()
	data, err := database.EmployeeDetailsRedisClient.Get(ctx, getEmployeeCacheKey(id)).Result()
	if err != nil {
		return nil, err
	}

	var employee models.EmployeeDetails
	if err := json.Unmarshal([]byte(data), &employee); err != nil {
		return nil, err
	}
	return &employee, nil
}

// FindEmployeeByID fetches an employee from MySQL by ID
func FindEmployeeByID(id uint) (*models.EmployeeDetails, error) {
	var employee models.EmployeeDetails

	// Query database using GORM
	result := database.EmployeeDetails.First(&employee, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &employee, nil
}

// below function will updated record in both mysql and redis
func UpdateEmployee(employee models.EmployeeDetails) error {
	// Update in MySQL
	result := database.EmployeeDetails.Save(&employee)
	if result.Error != nil {
		return result.Error
	}

	// Update in Redis
	err := CacheEmployee(employee)
	if err != nil {
		return err
	}

	return nil
}

// below function will delete record from both mysql and redis
func DeleteEmployee(id uint) error {
	// Delete from MySQL
	result := database.EmployeeDetails.Delete(&models.EmployeeDetails{}, id)
	if result.Error != nil {
		return result.Error
	}

	// Delete from Redis
	ctx := context.Background()
	err := database.EmployeeDetailsRedisClient.Del(ctx, getEmployeeCacheKey(id)).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetAllEmployees fetches all employees from MySQL
func GetAllEmployees() ([]models.EmployeeDetails, error) {
	var employees []models.EmployeeDetails

	// Query database using GORM
	result := database.EmployeeDetails.Find(&employees)
	if result.Error != nil {
		return nil, result.Error
	}

	return employees, nil
}
