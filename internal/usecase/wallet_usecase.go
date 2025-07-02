package usecase

import (
	"errors"
	"gameintegrationapi/internal/domain"
	"gameintegrationapi/internal/infrastructure"
	"gameintegrationapi/internal/repository"
	"log"
	"strconv"
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
	walletClient    *infrastructure.WalletClient
}

var ErrWalletServiceUnavailable = errors.New("wallet service is not available")

func NewWalletUseCase(userRepo repository.UserRepository, transactionRepo repository.TransactionRepository, db *gorm.DB, walletClient *infrastructure.WalletClient) WalletUseCase {
	return &walletUseCase{userRepo, transactionRepo, db, walletClient}
}

func (uc *walletUseCase) Withdraw(userID uint, amount float64, currency, providerTxID, roundID, gameID string) (*domain.Transaction, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("Withdraw: failed to find user: %v", err)
		return nil, err
	}
	walletID, err := strconv.ParseInt(user.WalletID, 10, 64)
	if err != nil {
		log.Printf("Withdraw: invalid wallet ID: %v", err)
		return nil, err
	}
	withdrawReq := infrastructure.WalletWithdrawRequest{
		Currency: user.Currency,
		Transactions: []struct {
			Amount    float64 `json:"amount"`
			BetID     int     `json:"betId"`
			Reference string  `json:"reference"`
		}{
			{
				Amount:    amount,
				BetID:     0,
				Reference: providerTxID,
			},
		},
		UserID: walletID,
	}
	_, err = uc.walletClient.Withdraw(withdrawReq)
	if err != nil {
		log.Printf("Withdraw: external wallet error: %v", err)
		return nil, err
	}
	if user.Balance < amount {
		log.Printf("Withdraw: insufficient funds for user %d", userID)
		return nil, errors.New("insufficient funds")
	}
	oldBalance := user.Balance
	newBalance := oldBalance - amount
	tx := &domain.Transaction{
		UserID:           userID,
		Type:             "WITHDRAW",
		Amount:           amount,
		OldBalance:       oldBalance,
		NewBalance:       newBalance,
		Status:           "COMPLETED",
		ProviderTxID:     providerTxID,
		ProviderRoundID:  roundID,
		ProviderGameID:   gameID,
		PlatformResponse: "{}",
		CreatedAt:        time.Now(),
	}
	err = uc.db.Transaction(func(txDb *gorm.DB) error {
		if err := repository.NewTransactionRepository(txDb).Create(tx); err != nil {
			log.Printf("Withdraw: failed to create transaction: %v", err)
			return err
		}
		if err := repository.NewUserRepository(txDb).UpdateBalance(user, newBalance); err != nil {
			log.Printf("Withdraw: failed to update balance: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("Withdraw: db transaction error: %v", err)
		return nil, err
	}
	log.Printf("Withdraw: success for user %d, amount %.2f", userID, amount)
	return tx, nil
}

func (uc *walletUseCase) Deposit(userID uint, amount float64, currency, providerTxID, providerParentTxID string) (*domain.Transaction, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("Deposit: failed to find user: %v", err)
		return nil, err
	}
	walletID, err := strconv.ParseInt(user.WalletID, 10, 64)
	if err != nil {
		log.Printf("Deposit: invalid wallet ID: %v", err)
		return nil, err
	}
	depositReq := infrastructure.WalletDepositRequest{
		Currency: user.Currency,
		Transactions: []struct {
			Amount    float64 `json:"amount"`
			BetID     int     `json:"betId"`
			Reference string  `json:"reference"`
		}{
			{
				Amount:    amount,
				BetID:     0,
				Reference: providerTxID,
			},
		},
		UserID: walletID,
	}
	_, err = uc.walletClient.Deposit(depositReq)
	if err != nil {
		log.Printf("Deposit: external wallet error: %v", err)
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
		PlatformResponse:   "{}",
		CreatedAt:          time.Now(),
	}
	err = uc.db.Transaction(func(txDb *gorm.DB) error {
		if err := repository.NewTransactionRepository(txDb).Create(tx); err != nil {
			log.Printf("Deposit: failed to create transaction: %v", err)
			return err
		}
		if err := repository.NewUserRepository(txDb).UpdateBalance(user, newBalance); err != nil {
			log.Printf("Deposit: failed to update balance: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("Deposit: db transaction error: %v", err)
		return nil, err
	}
	log.Printf("Deposit: success for user %d, amount %.2f", userID, amount)
	return tx, nil
}

func (uc *walletUseCase) Cancel(userID uint, providerTxID string) (*domain.Transaction, error) {
	originalTx, err := uc.transactionRepo.FindByProviderTxID(providerTxID)
	if err != nil {
		log.Printf("Cancel: original transaction not found: %v", err)
		return nil, errors.New("original transaction not found")
	}
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("Cancel: failed to find user: %v", err)
		return nil, err
	}
	walletID, err := strconv.ParseInt(user.WalletID, 10, 64)
	if err != nil {
		log.Printf("Cancel: invalid wallet ID: %v", err)
		return nil, err
	}
	cancelAmount := originalTx.Amount
	cancelReq := infrastructure.WalletDepositRequest{
		Currency: user.Currency,
		Transactions: []struct {
			Amount    float64 `json:"amount"`
			BetID     int     `json:"betId"`
			Reference string  `json:"reference"`
		}{
			{
				Amount:    cancelAmount,
				BetID:     0,
				Reference: "cancel-" + providerTxID,
			},
		},
		UserID: walletID,
	}
	_, err = uc.walletClient.Deposit(cancelReq)
	if err != nil {
		log.Printf("Cancel: external wallet error: %v", err)
		return nil, err
	}
	if originalTx.UserID != userID {
		log.Printf("Cancel: transaction does not belong to user %d", userID)
		return nil, errors.New("transaction does not belong to user")
	}
	if originalTx.Status == "CANCELLED" {
		log.Printf("Cancel: transaction already cancelled for user %d", userID)
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
		PlatformResponse:   "{}",
		CreatedAt:          time.Now(),
	}
	err = uc.db.Transaction(func(txDb *gorm.DB) error {
		if err := repository.NewTransactionRepository(txDb).Create(cancelTx); err != nil {
			log.Printf("Cancel: failed to create transaction: %v", err)
			return err
		}
		if err := repository.NewUserRepository(txDb).UpdateBalance(user, newBalance); err != nil {
			log.Printf("Cancel: failed to update balance: %v", err)
			return err
		}
		originalTx.Status = "CANCELLED"
		return txDb.Save(originalTx).Error
	})
	if err != nil {
		log.Printf("Cancel: db transaction error: %v", err)
		return nil, err
	}
	log.Printf("Cancel: success for user %d, amount %.2f", userID, cancelAmount)
	return cancelTx, nil
}
