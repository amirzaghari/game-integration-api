package usecase

import (
	"gameintegrationapi/internal/domain"
	"gameintegrationapi/internal/repository"
)

type PlayerUseCase interface {
	GetPlayerInfo(userID uint) (*domain.User, error)
}

type playerUseCase struct {
	userRepo repository.UserRepository
}

func NewPlayerUseCase(userRepo repository.UserRepository) PlayerUseCase {
	return &playerUseCase{userRepo}
}

func (uc *playerUseCase) GetPlayerInfo(userID uint) (*domain.User, error) {
	return uc.userRepo.FindByID(userID)
}
