package models

type UserModel struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size(100)"`
	Email    string `gorm:"unique"`
	Password string `gorm:"not null"`
	Status   string `gorm:"type:VARCHAR(20);check:status IN ('Active','Blocked')"`
}
