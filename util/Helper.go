package util

import "strings"

/*
* @Description : This Function will Clean String Data for SQL Injection Content
* @param data : The Data which will be cleaned
* @return : The Cleaned Data
 */
func CleanString(data string) string {
	// Remove any SQL injection patterns
	data = strings.ReplaceAll(data, "xp_", "")     // Remove SQL Server extended stored procedure prefix
	data = strings.ReplaceAll(data, "0x", "")      // Remove hexadecimal prefix
	data = strings.ReplaceAll(data, "0X", "")      // Remove hexadecimal prefix
	data = strings.ReplaceAll(data, "exec", "")    // Remove EXEC command
	data = strings.ReplaceAll(data, "execute", "") // Remove EXECUTE command

	return data
}
