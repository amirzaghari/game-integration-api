package repository

import (
	"gameintegrationapi/internal/domain"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *domain.Transaction) error
	FindByProviderTxID(providerTxID string) (*domain.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Create(tx *domain.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) FindByProviderTxID(providerTxID string) (*domain.Transaction, error) {
	var tx domain.Transaction
	if err := r.db.Where("provider_tx_id = ?", providerTxID).First(&tx).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}
