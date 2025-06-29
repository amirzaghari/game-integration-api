package infrastructure

import (
	"gameintegrationapi/internal/domain"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedTestUsers(db *gorm.DB) {
	users := []domain.User{
		{
			WalletID:  "34633089486",
			Username:  "testuser1",
			Password:  hashPassword("testpass"),
			Currency:  "USD",
			Balance:   5000.00,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			WalletID:  "34679664254",
			Username:  "testuser2",
			Password:  hashPassword("testpass"),
			Currency:  "EUR",
			Balance:   9000000000.00,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			WalletID:  "34616761765",
			Username:  "testuser3",
			Password:  hashPassword("testpass"),
			Currency:  "KES",
			Balance:   750.50,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			WalletID:  "34673635133",
			Username:  "testuser4",
			Password:  hashPassword("testpass"),
			Currency:  "USD",
			Balance:   31415.25,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	for _, u := range users {
		db.Where(domain.User{Username: u.Username}).FirstOrCreate(&u)
	}
}

func hashPassword(pw string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash)
}
