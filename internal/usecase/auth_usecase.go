package usecase

import (
	"errors"
	"gameintegrationapi/internal/repository"

	"gameintegrationapi/internal/infrastructure"
)

var jwtKey = []byte("your-secret-key") // TODO: Move to config

type AuthUseCase interface {
	Login(username, password string) (string, error)
}

type authUseCase struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(userRepo repository.UserRepository) AuthUseCase {
	return &authUseCase{userRepo}
}

func (uc *authUseCase) Login(username, password string) (string, error) {
	user, err := uc.userRepo.FindByCredentials(username, password)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := infrastructure.GenerateJWT(user.ID, user.Username, string(jwtKey))
	if err != nil {
		return "", err
	}

	return token, nil
}
