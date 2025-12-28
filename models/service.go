package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(dbPath string) (*Service, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
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
		if err := s.DB.Create(&b).Error; err != nil {
			return err
		}
	}

	// Seed Accounts
	accounts := []Account{
		{
			Name:          "CapitalOne",
			Description:   "Capital One rewards Credit Account",
			BeneficiaryID: "Us",
			// Note: We haven't added CSV config fields to Account struct in the final models.go!
			// Wait, the plan said we would, but I reverted strict fields and decided to use hardcoded parsers mapped by Name?
			// Let's check the Account struct again. It DOES NOT have CSVFormat.
			// The Rule said: "The CSV format is different for each account... See table below".
			// It implies we know the mapping.
			// I will use the Account Name to determine parsing logic in the backend, as the simplest approach.
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
		if err := s.DB.Create(&a).Error; err != nil {
			return err
		}
	}

	return nil
}
