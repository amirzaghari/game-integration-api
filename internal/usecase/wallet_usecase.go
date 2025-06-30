package usecase

import (
	"errors"
	"gameintegrationapi/internal/domain"
	"gameintegrationapi/internal/repository"
	"time"

	"gorm.io/gorm"
)

type WalletUseCase interface {
	Withdraw(userID uint, amount float64, currency, providerTxID, roundID, gameID string) (*domain.Transaction, error)
	Deposit(userID uint, amount float64, currency, providerTxID, providerParentTxID string) (*domain.Transaction, error)
	Cancel(userID uint, providerTxID string) (*domain.Transaction, error)
}

type walletUseCase struct {
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
	db              *gorm.DB
}

func NewWalletUseCase(userRepo repository.UserRepository, transactionRepo repository.TransactionRepository, db *gorm.DB) WalletUseCase {
	return &walletUseCase{userRepo, transactionRepo, db}
}

func (uc *walletUseCase) Withdraw(userID uint, amount float64, currency, providerTxID, roundID, gameID string) (*domain.Transaction, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user.Balance < amount {
		return nil, errors.New("insufficient funds")
	}

	oldBalance := user.Balance
	newBalance := oldBalance - amount

	tx := &domain.Transaction{
		UserID:          userID,
		Type:            "WITHDRAW",
		Amount:          amount,
		OldBalance:      oldBalance,
		NewBalance:      newBalance,
		Status:          "COMPLETED",
		ProviderTxID:    providerTxID,
		ProviderRoundID: roundID,
		ProviderGameID:  gameID,
		CreatedAt:       time.Now(),
	}

	err = uc.db.Transaction(func(txDb *gorm.DB) error {
		if err := repository.NewTransactionRepository(txDb).Create(tx); err != nil {
			return err
		}
		if err := repository.NewUserRepository(txDb).UpdateBalance(user, newBalance); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (uc *walletUseCase) Deposit(userID uint, amount float64, currency, providerTxID, providerParentTxID string) (*domain.Transaction, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	oldBalance := user.Balance
	newBalance := oldBalance + amount
	status := "WON"
	if amount == 0 {
		status = "LOST"
	}

	tx := &domain.Transaction{
		UserID:             userID,
		Type:               "DEPOSIT",
		Amount:             amount,
		OldBalance:         oldBalance,
		NewBalance:         newBalance,
		Status:             status,
		ProviderTxID:       providerTxID,
		ProviderParentTxID: providerParentTxID,
		CreatedAt:          time.Now(),
	}

	err = uc.db.Transaction(func(txDb *gorm.DB) error {
		if err := repository.NewTransactionRepository(txDb).Create(tx); err != nil {
			return err
		}
		if err := repository.NewUserRepository(txDb).UpdateBalance(user, newBalance); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (uc *walletUseCase) Cancel(userID uint, providerTxID string) (*domain.Transaction, error) {
	originalTx, err := uc.transactionRepo.FindByProviderTxID(providerTxID)
	if err != nil {
		return nil, errors.New("original transaction not found")
	}

	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if originalTx.UserID != userID {
		return nil, errors.New("transaction does not belong to user")
	}

	if originalTx.Status == "CANCELLED" {
		return nil, errors.New("transaction already cancelled")
	}

	oldBalance := user.Balance
	newBalance := oldBalance + originalTx.Amount

	cancelTx := &domain.Transaction{
		UserID:             userID,
		Type:               "CANCEL",
		Amount:             originalTx.Amount,
		OldBalance:         oldBalance,
		NewBalance:         newBalance,
		Status:             "CANCELLED",
		ProviderTxID:       "cancel-" + providerTxID,
		ProviderParentTxID: providerTxID,
		CreatedAt:          time.Now(),
	}

	err = uc.db.Transaction(func(txDb *gorm.DB) error {
		if err := repository.NewTransactionRepository(txDb).Create(cancelTx); err != nil {
			return err
		}
		if err := repository.NewUserRepository(txDb).UpdateBalance(user, newBalance); err != nil {
			return err
		}
		originalTx.Status = "CANCELLED"
		return txDb.Save(originalTx).Error
	})

	if err != nil {
		return nil, err
	}

	return cancelTx, nil
}
