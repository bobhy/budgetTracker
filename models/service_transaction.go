package models

import (
	"fmt"
	"math/rand"
	"time"
)

func (s *Service) GetTransactions() ([]Transaction, error) {
	var transactions []Transaction
	// Preload Account?
	result := s.db.Preload("Account").Find(&transactions)
	return transactions, result.Error
}

func (s *Service) AddTransaction(postedDate Date, accountID string, amount Money, description, tag string) error {
	return s.db.Create(&Transaction{
		PostedDate:  postedDate,
		AccountID:   accountID,
		Amount:      amount,
		Description: description,
		Tag:         tag,
	}).Error
}

// UpdateTransaction uses ID for identification
func (s *Service) UpdateTransaction(id uint, postedDate Date, accountID string, amount Money, description, tag string) error {
	return s.db.Model(&Transaction{}).Where("id = ?", id).Updates(map[string]interface{}{
		"posted_date": postedDate,
		"account_id":  accountID,
		"amount":      amount,
		"description": description,
		"tag":         tag,
	}).Error
}

func (s *Service) DeleteTransaction(id uint) error {
	return s.db.Delete(&Transaction{}, id).Error
}

func (s *Service) GenerateTransactions(count int) error {
	accs, err := s.GetAccounts()
	if err != nil {
		return err
	}
	if len(accs) == 0 {
		return fmt.Errorf("no accounts found")
	}

	for i := 0; i < count; i++ {
		acc := accs[rand.Intn(len(accs))]
		date := time.Now().AddDate(0, 0, -rand.Intn(365)).Format("2006-01-02")
		amount := Money(rand.Int63n(20000) - 10000) // -100.00 to +100.00
		if err := s.AddTransaction(Date(date), acc.Name, amount, "Random transaction", "random"); err != nil {
			return err
		}
	}
	return nil
}
