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

// Below Function will get all employees from mysql
// and return paginated results
// Pagination is done using limit and offset
func GetAllEmployees(page int, limit int) ([]models.EmployeeDetails, error) {
	var employees []models.EmployeeDetails

	// Calculate offset
	offset := (page - 1) * limit

	// Query database with limit and offset
	result := database.EmployeeDetails.
		Limit(limit).
		Offset(offset).
		Find(&employees)

	if result.Error != nil {
		return nil, result.Error
	}

	return employees, nil
}

// SyncAndClearEmployeeCache syncs Redis data to MySQL and clears Redis cache
func SyncAndClearEmployeeCache() error {
	ctx := context.Background()
	var cursor uint64
	var keys []string
	var err error

	// Scan for all keys matching "employee:*"
	for {
		var tempKeys []string
		tempKeys, cursor, err = database.EmployeeDetailsRedisClient.Scan(ctx, cursor, "employee:*", 100).Result()
		if err != nil {
			return fmt.Errorf("error scanning redis keys: %v", err)
		}
		keys = append(keys, tempKeys...)
		if cursor == 0 {
			break
		}
	}

	// Loop through all keys
	for _, key := range keys {
		// Get JSON string from Redis
		data, err := database.EmployeeDetailsRedisClient.Get(ctx, key).Result()
		if err != nil {
			fmt.Printf("Failed to get key %s: %v\n", key, err)
			continue
		}

		var employee models.EmployeeDetails
		if err := json.Unmarshal([]byte(data), &employee); err != nil {
			fmt.Printf("Failed to unmarshal data for key %s: %v\n", key, err)
			continue
		}

		// Update MySQL (this also re-caches, which we will delete after)
		if err := UpdateEmployee(employee); err != nil {
			fmt.Printf("Failed to update employee %d in DB: %v\n", employee.ID, err)
			continue
		}

		// Delete from Redis cache
		if err := database.EmployeeDetailsRedisClient.Del(ctx, key).Err(); err != nil {
			fmt.Printf("Failed to delete Redis key %s: %v\n", key, err)
		}
	}

	return nil
}
