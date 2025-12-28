package models

func (s *Service) GetAccounts() ([]Account, error) {
	var accounts []Account
	// Preload implicit if needed, but for simple list, just data.
	// Maybe preload Beneficiary just in case
	result := s.DB.Find(&accounts)
	return accounts, result.Error
}

func (s *Service) AddAccount(name, description, beneficiaryID string) error {
	return s.DB.Create(&Account{
		Name:          name,
		Description:   description,
		BeneficiaryID: beneficiaryID,
	}).Error
}

func (s *Service) UpdateAccount(oldName, newName, description, beneficiaryID string) error {
	// If PK (Name) is changing
	if oldName != newName {
		return s.DB.Model(&Account{Name: oldName}).Updates(Account{
			Name:          newName,
			Description:   description,
			BeneficiaryID: beneficiaryID,
		}).Error
	}
	// Normal update
	return s.DB.Model(&Account{Name: oldName}).Updates(Account{
		Description:   description,
		BeneficiaryID: beneficiaryID,
	}).Error
}

func (s *Service) DeleteAccount(name string) error {
	return s.DB.Delete(&Account{Name: name}).Error
}
