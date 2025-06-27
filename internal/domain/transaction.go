package domain

import "time"

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
