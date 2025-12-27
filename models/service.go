package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(dbPath string) (*Service, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migrate
	err = db.AutoMigrate(&Beneficiary{}, &Account{}, &Budget{}, &Transaction{})
	if err != nil {
		return nil, err
	}

	return &Service{db: db}, nil
}

// Clean drops all tables and re-migrates them
func (s *Service) Clean() error {
	// Drop tables
	err := s.db.Migrator().DropTable(&Beneficiary{}, &Account{}, &Budget{}, &Transaction{})
	if err != nil {
		return err
	}
	// Re-migrate
	return s.db.AutoMigrate(&Beneficiary{}, &Account{}, &Budget{}, &Transaction{})
}
