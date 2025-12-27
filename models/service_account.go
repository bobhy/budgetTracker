package models

import (
	"fmt"
	"math/rand"
)

func (s *Service) GetAccounts() ([]Account, error) {
	var accounts []Account
	// Preload implicit if needed, but for simple list, just data.
	// Maybe preload Beneficiary just in case
	result := s.db.Find(&accounts)
	return accounts, result.Error
}

func (s *Service) AddAccount(name, description, beneficiaryID string) error {
	return s.db.Create(&Account{
		Name:          name,
		Description:   description,
		BeneficiaryID: beneficiaryID,
	}).Error
}

func (s *Service) UpdateAccount(oldName, newName, description, beneficiaryID string) error {
	// If PK (Name) is changing
	if oldName != newName {
		return s.db.Model(&Account{Name: oldName}).Updates(Account{
			Name:          newName,
			Description:   description,
			BeneficiaryID: beneficiaryID,
		}).Error
	}
	// Normal update
	return s.db.Model(&Account{Name: oldName}).Updates(Account{
		Description:   description,
		BeneficiaryID: beneficiaryID,
	}).Error
}

func (s *Service) DeleteAccount(name string) error {
	return s.db.Delete(&Account{Name: name}).Error
}

func (s *Service) GenerateAccounts(count int) error {
	bens, err := s.GetBeneficiaries()
	if err != nil {
		return err
	}
	if len(bens) == 0 {
		return fmt.Errorf("no beneficiaries to link accounts to")
	}

	for i := 0; i < count; i++ {
		ben := bens[rand.Intn(len(bens))]
		name := fmt.Sprintf("Account %s", randomString(5))
		if err := s.AddAccount(name, "Randomly generated account", ben.Name); err != nil {
			return err
		}
	}
	return nil
}
