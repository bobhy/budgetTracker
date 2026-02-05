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
}

type Account struct {
	Name           string `gorm:"primaryKey;default:'None';constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT"`
	Description    string
	Beneficiary    string
	BeneficiaryObj *Beneficiary `gorm:"foreignKey:Beneficiary;references:Name"`
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
	BeneficiaryObj *Beneficiary `gorm:"foreignKey:Beneficiary;references:Name"`
	Amount         Money
	IntervalMonths int
}

type Transaction struct {
	ID             uint       `gorm:"primarykey;autoIncrement"`
	CreatedAt      time.Time  `json:"-"` // Hide from frontend to avoid warning/binding issues
	UpdatedAt      time.Time  `json:"-"`
	DeletedAt      *time.Time `gorm:"index"` // Use pointer to time for soft delete, hide from json
	PostedDate     Date
	Account        string
	AccountObj     *Account `gorm:"foreignKey:Account;references:Name"`
	Amount         Money
	Description    string
	Tag            string
	Budget         string
	BudgetObj      *Budget `gorm:"foreignKey:Budget;references:Name"`
	Beneficiary    string
	BeneficiaryObj *Beneficiary `gorm:"foreignKey:Beneficiary;references:Name"` // Overrides Account default if set
	RawHint        string       // Category hint from import
}

// RawTransaction is used for importing transactions before they are fully processed and linked
type RawTransaction struct {
	ID             uint       `gorm:"primarykey;autoIncrement"`
	CreatedAt      time.Time  `json:"-"`
	UpdatedAt      time.Time  `json:"-"`
	DeletedAt      *time.Time `gorm:"index"`
	PostedDate     Date
	Account        string
	AccountObj     *Account `gorm:"foreignKey:Account;references:Name"`
	Amount         Money
	Description    string
	Tag            string
	Budget         string
	BudgetObj      *Budget `gorm:"foreignKey:Budget;references:Name"`
	Action         string  // "add" or "update"
	Beneficiary    string
	BeneficiaryObj *Beneficiary `gorm:"foreignKey:Beneficiary;references:Name"`
	RawHint        string
}
