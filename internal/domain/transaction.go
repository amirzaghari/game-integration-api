package domain

import "time"

type Transaction struct {
	ID                 uint    `gorm:"primaryKey"`
	UserID             uint    `gorm:"index;not null"`
	BetID              uint    `gorm:"index"`
	Type               string  `gorm:"not null"` // WITHDRAW, DEPOSIT, CANCEL
	Amount             float64 `gorm:"not null"`
	OldBalance         float64 `gorm:"not null"`
	NewBalance         float64 `gorm:"not null"`
	Status             string  `gorm:"not null"` // WON, LOST, CANCELLED
	ProviderTxID       string  `gorm:"index"`
	ProviderParentTxID string  `gorm:"index"` // To link deposit/cancel to original withdraw
	ProviderRoundID    string  `gorm:"index"`
	ProviderGameID     string  `gorm:"index"`
	ProviderSessionID  string  `gorm:"index"`
	PlatformResponse   string  `gorm:"type:jsonb"`
	CreatedAt          time.Time
}
