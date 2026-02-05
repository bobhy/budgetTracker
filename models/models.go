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
	Name string `gorm:"primaryKey;default:'None';constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT"`

	// Relationships
	Accounts        []Account        `gorm:"foreignKey:Beneficiary;references:Name"`
	Budgets         []Budget         `gorm:"foreignKey:Beneficiary;references:Name"`
	Transactions    []Transaction    `gorm:"foreignKey:Beneficiary;references:Name"`
	RawTransactions []RawTransaction `gorm:"foreignKey:Beneficiary;references:Name"`
}

type Account struct {
	Name        string `gorm:"primaryKey;default:'None';constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT"`
	Description string
	Beneficiary string

	// Relationships
	Budgets         []Budget         `gorm:"foreignKey:Beneficiary;references:Name"`
	Transactions    []Transaction    `gorm:"foreignKey:Account;references:Name"`
	RawTransactions []RawTransaction `gorm:"foreignKey:Account;references:Name"`
}

// Budget represents a planned expenditure over time
// It is possible for each beneficiary to have a budget
// for the same category of spending.
// It's also possible there is only one budget with a given name
// and for that budget to be assigned to *any* beneficiary.
//
// DB-wise, name and beneficiary should be a composite primary key
// but that's awkward to use when assigning budgets to transactions.
// Instead, we'll use a simple naming convention and enforce it
// When updating the budget table.
//
// 1. name shall be a unique primary key.  beneficiary
// shall not be part of the db-enforced constraint.
// 2. if you want to add a second budget item with the same name
// but for a different beneficiary, call it `<name>_<beneficiary>`.
// So Bob can have his boat {"boat", "bob"},
// Jessie can have her art {"art", "jessie"},
// but there can be shared and personal travel accounts:
// {"travel", "us"},
// {"travel_bob", "bob"},
// {"travel_jessie", "jessie"}
// 3. don't use '_' in name except for this distinction.
type Budget struct {
	Name           string `gorm:"primaryKey;default:'None';constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT"`
	Description    string
	Beneficiary    string
	Amount         Money
	IntervalMonths int

	// Relationships
	Transactions    []Transaction    `gorm:"foreignKey:Budget;references:Name"`
	RawTransactions []RawTransaction `gorm:"foreignKey:Budget;references:Name"`
}

type Transaction struct {
	ID          uint       `gorm:"primarykey"`
	CreatedAt   time.Time  `json:"-"` // Hide from frontend to avoid warning/binding issues
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `gorm:"index"` // Use pointer to time for soft delete, hide from json
	PostedDate  Date
	Account     string `gorm:"foreignKey:Account"`
	Amount      Money
	Description string
	Tag         string
	Budget      string
	Beneficiary string // Overrides Account default if set
	RawHint     string // Category hint from import
}

// RawTransaction is used for importing transactions before they are fully processed and linked
type RawTransaction struct {
	ID          uint       `gorm:"primarykey"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `gorm:"index"`
	PostedDate  Date
	Account     string `gorm:"foreignKey:Account"`
	Amount      Money
	Description string
	Tag         string
	Budget      string `gorm:"foreignKey:Budget"`
	Action      string // "add" or "update"
	Beneficiary string `gorm:"foreignKey:Beneficiary"`
	RawHint     string
}
