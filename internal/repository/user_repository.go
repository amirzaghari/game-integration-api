package repository

import (
	"errors"
	"gameintegrationapi/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByCredentials(username, password string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	UpdateBalance(user *domain.User, newBalance float64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByCredentials(username, password string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateBalance(user *domain.User, newBalance float64) error {
	return r.db.Model(user).Update("balance", newBalance).Error
}
