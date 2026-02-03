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
	Name string `gorm:"primaryKey" json:"name"`

	// Relationships
	Accounts []Account `gorm:"foreignKey:BeneficiaryID" json:"accounts,omitempty"`
	Budgets  []Budget  `gorm:"foreignKey:BeneficiaryID" json:"budgets,omitempty"`
}

type Account struct {
	Name          string `gorm:"primaryKey" json:"name"`
	Description   string `json:"description"`
	BeneficiaryID string `json:"beneficiary_id"`

	// Relationships
	Transactions []Transaction `gorm:"foreignKey:AccountID" json:"transactions,omitempty"`
}

type Budget struct {
	Name           string `gorm:"primaryKey" json:"name"`
	Description    string `json:"description"`
	BeneficiaryID  string `json:"beneficiary_id"`
	Amount         Money  `json:"amount"`
	IntervalMonths int    `json:"interval_months"`
}

type Transaction struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"-"` // Hide from frontend to avoid warning/binding issues
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `gorm:"index" json:"-"` // Use pointer to time for soft delete, hide from json
	PostedDate  Date       `json:"posted_date"`
	AccountID   string     `json:"account_id"`
	Account     Account    `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Amount      Money      `json:"amount"`
	Description string     `json:"description"`
	Tag         string     `json:"tag"`
	BudgetLine  string     `json:"budget_line"` // Replaces BudgetLineID
	Beneficiary string     `json:"beneficiary"` // Overrides Account default if set
	RawHint     string     `json:"raw_hint"`    // Category hint from import
}

// RawTransaction is used for importing transactions before they are fully processed and linked
type RawTransaction struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `gorm:"index" json:"-"`
	PostedDate  Date       `json:"posted_date"`
	AccountID   string     `json:"account_id"`
	Amount      Money      `json:"amount"`
	Description string     `json:"description"`
	Tag         string     `json:"tag"`
	BudgetLine  string     `json:"budget_line"`
	Action      string     `json:"action"` // "add" or "update"
	Beneficiary string     `json:"beneficiary"`
	RawHint     string     `json:"raw_hint"`
}
