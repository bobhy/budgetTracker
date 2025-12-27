package models

import (
	"gorm.io/gorm"
)

// Money is stored as integer number of cents
type Money int64

// Date is stored as string YYYY-MM-DD
type Date string

type Beneficiary struct {
	Name string `gorm:"primaryKey"`

	// Relationships
	Accounts []Account `gorm:"foreignKey:BeneficiaryID"`
	Budgets  []Budget  `gorm:"foreignKey:BeneficiaryID"`
}

type Account struct {
	Name          string `gorm:"primaryKey"`
	Description   string
	BeneficiaryID string

	// Relationships
	Transactions []Transaction `gorm:"foreignKey:AccountID"`
}

type Budget struct {
	Name           string `gorm:"primaryKey"`
	Description    string
	BeneficiaryID  string
	Amount         Money
	IntervalMonths int
}

type Transaction struct {
	gorm.Model
	PostedDate   Date
	AccountID    string
	Account      Account `gorm:"foreignKey:AccountID"`
	Amount       Money
	Description  string
	Tag          string
	BudgetLineID *uint // Optional link to a budget line item (not fully defined in design yet, using simple ID for now)
}
