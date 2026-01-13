package models

import (
	"time"
)

// Money is stored as integer number of cents
type Money int64

// Date is stored as string YYYY-MM-DD
type Date string

type SortOption struct {
	Key       string `json:"key"`
	Direction string `json:"direction"` // "asc" or "desc"
}

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
	ID          uint       `gorm:"primarykey" json:"ID"`
	CreatedAt   time.Time  `json:"-"` // Hide from frontend to avoid warning/binding issues
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `gorm:"index" json:"-"` // Use pointer to time for soft delete, hide from json
	PostedDate  Date
	AccountID   string
	Account     Account `gorm:"foreignKey:AccountID"`
	Amount      Money
	Description string
	Tag         string
	BudgetLine  string // Replaces BudgetLineID
	Beneficiary string // Overrides Account default if set
	RawHint     string // Category hint from import
}

// RawTransaction is used for importing transactions before they are fully processed and linked
type RawTransaction struct {
	ID          uint       `gorm:"primarykey" json:"ID"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `gorm:"index" json:"-"`
	PostedDate  Date
	AccountID   string
	Amount      Money
	Description string
	Tag         string
	BudgetLine  string
	Action      string // "add" or "update"
	Beneficiary string
	RawHint     string
}
