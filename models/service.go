package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(dbPath string) (*Service, error) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s?%s", dbPath, "_foreign_keys=on")), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migrate
	err = db.AutoMigrate(&Beneficiary{}, &Account{}, &Budget{}, &Transaction{}, &RawTransaction{})
	if err != nil {
		return nil, err
	}

	return &Service{DB: db}, nil
}

// Clean drops all tables and re-migrates them, then seeds production data
func (s *Service) Clean() error {
	// Drop tables
	err := s.DB.Migrator().DropTable(&Beneficiary{}, &Account{}, &Budget{}, &Transaction{}, &RawTransaction{})
	if err != nil {
		return err
	}
	// Re-migrate
	err = s.DB.AutoMigrate(&Beneficiary{}, &Account{}, &Budget{}, &Transaction{}, &RawTransaction{})
	if err != nil {
		return err
	}

	return s.seed()
}

func (s *Service) seed() error {
	// Seed Beneficiaries
	beneficiaries := []Beneficiary{
		{Name: "Bob"},
		{Name: "Jessie"},
		{Name: "Us"},
	}
	for _, b := range beneficiaries {
		if err := Create(s.DB, &b); err != nil {
			return err
		}
	}

	// Seed Accounts
	accounts := []Account{
		{
			Name:          "CapitalOne",
			Description:   "Capital One rewards Credit Account",
			BeneficiaryID: "Us",
		},
		{
			Name:          "WfChecking",
			Description:   "Wells Fargo checking",
			BeneficiaryID: "Us",
		},
		{
			Name:          "WfVisa",
			Description:   "Wells Fargo Visa",
			BeneficiaryID: "Us",
		},
	}
	for _, a := range accounts {
		if err := Create(s.DB, &a); err != nil {
			return err
		}
	}

	return nil
}
