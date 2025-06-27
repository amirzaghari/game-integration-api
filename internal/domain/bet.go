package domain

import "time"

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
