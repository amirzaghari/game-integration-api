package domain

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	WalletID  string `gorm:"uniqueIndex;not null"`
	Username  string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Currency  string `gorm:"not null"`
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}
