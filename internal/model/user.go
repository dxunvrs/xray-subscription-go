package model

type User struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Email string `gorm:"size:255;not null;uniqueIndex"`
	UUID  string `gorm:"size:36;not null;uniqueIndex"`
}
