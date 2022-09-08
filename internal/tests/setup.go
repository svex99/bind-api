package tests

import (
	"testing"

	"github.com/svex99/bind-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDatabase(t *testing.T) {
	DB, err := gorm.Open(sqlite.Open("data/bind-api-test.db"), &gorm.Config{
		AllowGlobalUpdate: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := DB.Exec("PRAGMA foreign_keys=ON").Error; err != nil {
		t.Fatal(err)
	}

	if err := DB.AutoMigrate(
		&models.User{},
		&models.Domain{},
		&models.SOARecord{},
		&models.NSRecord{},
		&models.ARecord{},
		&models.MXRecord{},
		&models.TXTRecord{},
	); err != nil {
		t.Fatal(err)
	}

	models.DB = DB

	if err := DB.Delete(&models.User{}).Error; err != nil {
		t.Fatal(err)
	}
	if err := DB.Delete(&models.Domain{}).Error; err != nil {
		t.Fatal(err)
	}
	if err := DB.Delete(&models.SOARecord{}).Error; err != nil {
		t.Fatal(err)
	}
	if err := DB.Delete(&models.NSRecord{}).Error; err != nil {
		t.Fatal(err)
	}
	if err := DB.Delete(&models.ARecord{}).Error; err != nil {
		t.Fatal(err)
	}
	if err := DB.Delete(&models.MXRecord{}).Error; err != nil {
		t.Fatal(err)
	}
	if err := DB.Delete(&models.TXTRecord{}).Error; err != nil {
		t.Fatal(err)
	}
}
