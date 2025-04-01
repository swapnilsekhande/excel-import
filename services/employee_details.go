package services

import (
	"errors"
	"excel-import/models"
	"excel-import/repositories"
	"excel-import/util"
	"fmt"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sirupsen/logrus"
)

// Get employee details (check Redis first, fallback to DB)
func GetEmployeeByID(id uint) (*models.EmployeeDetails, error) {
	// Try fetching from Redis
	employee, err := repositories.GetCachedEmployee(id)
	if err == nil {
		return employee, nil
	}

	// Fetch from database
	employee, err = repositories.FindEmployeeByID(id)
	if err != nil {
		return nil, err
	}

	// Store in Redis for future requests
	repositories.CacheEmployee(*employee)

	return employee, nil
}

/*
* @Description: This function is used to read the excel file and return the data in the form of a slice of EmployeeDetails
* @param filePath: The path of the excel file
* @return []EmployeeDetails: The slice of EmployeeDetails
* @return error: The error if any
 */
func ReadExcelFile(filePath string) ([]models.EmployeeDetails, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		logrus.Error("Error opening Excel file:", err)
		return nil, err
	}

	sheetName := f.GetSheetName(1)
	if sheetName == "" {
		logrus.Error("No sheets found in the Excel file")
		return nil, errors.New("no sheets found")
	}

	rows := f.GetRows(sheetName)
	if len(rows) == 0 {
		return nil, errors.New("empty Excel file")
	}

	if len(rows) == 0 {
		return nil, errors.New("empty Excel file")
	}

	headers := rows[0]
	requiredHeaders := []string{"first_name", "last_name", "company_name", "address", "city", "county", "postal", "phone", "email", "web"}
	for i, header := range requiredHeaders {
		if i >= len(headers) || strings.TrimSpace(headers[i]) != header {
			err := fmt.Errorf("missing or invalid header: %s", header)
			logrus.Error(err)
			return nil, err
		}
	}

	var employeeDetails []models.EmployeeDetails

	for _, row := range rows[1:] { // Skip header row
		if len(row) < 9 {
			continue // Skip rows with insufficient data
		}
		employee := models.EmployeeDetails{
			FirstName:   util.CleanString(strings.TrimSpace(row[0])),
			LastName:    util.CleanString(strings.TrimSpace(row[1])),
			CompanyName: util.CleanString(strings.TrimSpace(row[2])),
			Address:     util.CleanString(strings.TrimSpace(row[3])),
			City:        util.CleanString(strings.TrimSpace(row[4])),
			County:      util.CleanString(strings.TrimSpace(row[5])),
			Postal:      util.CleanString(strings.TrimSpace(row[6])),
			Phone:       util.CleanString(strings.TrimSpace(row[7])),
			Email:       util.CleanString(strings.TrimSpace(row[8])),
			Web:         util.CleanString(strings.TrimSpace(row[9])),
		}
		employeeDetails = append(employeeDetails, employee)
	}

	if len(employeeDetails) == 0 {
		logrus.Error("no valid employee data found")
		return nil, errors.New("no valid employee data found")
	}

	return employeeDetails, nil
}
