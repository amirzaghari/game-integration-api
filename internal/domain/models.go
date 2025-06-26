package domain

import (
	"time"
)

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

type Bet struct {
	ID            uint    `gorm:"primaryKey"`
	UserID        uint    `gorm:"index;not null"`
	ProviderTxID  string  `gorm:"uniqueIndex;not null"`
	Amount        float64 `gorm:"not null"`
	Status        string  `gorm:"not null"` // PLACED, WON, LOST, CANCELLED
	WithdrawnTxID string  // For linking deposit to withdrawal
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Transaction struct {
	ID           uint    `gorm:"primaryKey"`
	UserID       uint    `gorm:"index;not null"`
	BetID        uint    `gorm:"index"`
	Type         string  `gorm:"not null"` // WITHDRAW, DEPOSIT, CANCEL
	Amount       float64 `gorm:"not null"`
	OldBalance   float64 `gorm:"not null"`
	NewBalance   float64 `gorm:"not null"`
	Status       string  `gorm:"not null"`
	ProviderTxID string  `gorm:"index"`
	CreatedAt    time.Time
}
