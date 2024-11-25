package initializers

import "main.go/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.UserModel{},&models.AdminModel{})
}
