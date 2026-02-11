/**
* models.go
*
* Models for the database.
* Uses GORM --> WAILS to marshal between DB, Golang and JS front end.
*
* todo:
*	make string match case insensitive (needs `collate nocase` in schema, but GORM doesn't?)
 */
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

// Beneficiary is a person who owns an account or a budget line
type Beneficiary struct {
	Name string `gorm:"primaryKey;default:'None';constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT"`
}

// Account is a place where money is held.
// We merge funds from multiple accounts into a budget line
type Account struct {
	Name           string `gorm:"primaryKey;default:'None';constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT"`
	Description    string
	Beneficiary    string
	BeneficiaryObj *Beneficiary `gorm:"foreignKey:Beneficiary;references:Name" json:"-"`
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
	BeneficiaryObj *Beneficiary `gorm:"foreignKey:Beneficiary;references:Name" json:"-"`
	Amount         Money
	IntervalMonths int
}

// A financial event in an account
type Transaction struct {
	ID             uint       `gorm:"primarykey;autoIncrement"`
	CreatedAt      time.Time  `json:"-"` // Hide from frontend to avoid warning/binding issues
	UpdatedAt      time.Time  `json:"-"`
	DeletedAt      *time.Time `gorm:"index"` // Use pointer to time for soft delete, hide from json
	PostedDate     Date       `gorm:"column:posted_date"`
	Account        string
	AccountObj     *Account `gorm:"foreignKey:Account;references:Name" json:"-"`
	Amount         Money
	Description    string // Descriptive text as provided by the bank
	Budget         string
	BudgetObj      *Budget `gorm:"foreignKey:Budget;references:Name" json:"-"`
	Beneficiary    string
	BeneficiaryObj *Beneficiary `gorm:"foreignKey:Beneficiary;references:Name" json:"-"` // Overrides Account default if set
	RawHint        string       // Category hint from import
}

// RawTransaction is used for importing transactions before they are fully processed and linked
// No FK constraints here, to permit import of raw data
type RawTransaction struct {
	ID          uint       `gorm:"primarykey;autoIncrement"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `gorm:"index"`
	PostedDate  Date       `gorm:"column:posted_date"`
	Account     string
	Amount      Money
	Description string
	Tag         string // Tag, usually derived from the Description. *not* a foreign key
	Budget      string // *not* a foreign key so we can import garbage from CSV
	Action      string // "add" or "update"
	Beneficiary string
	RawHint     string
}

// Tag is a string mapped to a Budget
type Tag struct {
	Name      string `gorm:"primaryKey;default:'';constraint:OnUpdate:CASCADE,OnDelete:SET DEFAULT"`
	Budget    string
	BudgetObj *Budget `gorm:"foreignKey:Budget;references:Name" json:"-"`
}
