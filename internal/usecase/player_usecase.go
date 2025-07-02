package usecase

import (
	"gameintegrationapi/internal/domain"
	"gameintegrationapi/internal/infrastructure"
	"gameintegrationapi/internal/repository"
	"log"
	"strconv"
)

type PlayerUseCase interface {
	GetPlayerInfo(userID uint) (*domain.User, error)
}

type playerUseCase struct {
	userRepo     repository.UserRepository
	walletClient *infrastructure.WalletClient
}

func NewPlayerUseCase(userRepo repository.UserRepository, walletClient *infrastructure.WalletClient) PlayerUseCase {
	return &playerUseCase{userRepo, walletClient}
}

func (uc *playerUseCase) GetPlayerInfo(userID uint) (*domain.User, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("GetPlayerInfo: failed to find user: %v", err)
		return nil, err
	}
	walletID, err := strconv.ParseInt(user.WalletID, 10, 64)
	if err != nil {
		log.Printf("GetPlayerInfo: invalid wallet ID: %v", err)
		return nil, err
	}
	profile, err := uc.walletClient.GetBalance(walletID)
	if err != nil {
		if err == infrastructure.ErrWalletUserNotFound {
			return nil, infrastructure.ErrWalletUserNotFound
		}
		log.Printf("GetPlayerInfo: external wallet error: %v", err)
		return nil, err
	}
	user.Balance = profile.Balance
	user.Currency = profile.Currency
	log.Printf("GetPlayerInfo: success for user %d", userID)
	return user, nil
}
