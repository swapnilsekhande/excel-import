package models

type EmployeeDetails struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	FirstName   string `json:"first_name" gorm:"not null" binding:"required,alpha"`
	LastName    string `json:"last_name" gorm:"not null" binding:"required,alpha"`
	CompanyName string `json:"company_name" gorm:"not null" binding:"required"`
	Address     string `json:"address" gorm:"not null" binding:"required"`
	City        string `json:"city" gorm:"not null" binding:"required,alpha"`
	County      string `json:"county" gorm:"not null" binding:"required,alpha"`
	Postal      string `json:"postal" gorm:"not null" binding:"required,alphanum"`
	Phone       string `json:"phone" gorm:"not null" binding:"required,numeric"`
	Email       string `json:"email" gorm:"not null" binding:"required,email"` // Should Make Unqiue but Repeated in Excel
	Web         string `json:"web" gorm:"not null" binding:"required,url"`
}
