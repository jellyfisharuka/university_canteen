package models



type Users struct {
	ID           int      `gorm:"type:int;primary_key"`
	UserName    string    `gorm:"type:varchar(255);not null"`
	Email       string    `gorm:"uniqueIndex;not null"`
	Password    string    `gorm:"not null"`
}