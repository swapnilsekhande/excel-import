package migrations

import (
	"excel-import/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func MigrationRun(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "2025030101",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().CreateTable(&models.EmployeeDetails{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("employee_details")
			},
		},
	})

	if err := m.Migrate(); err != nil {
		logrus.Error("Migration failed: ", err)
	}
	return nil
}
