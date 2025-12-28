package models

func (s *Service) GetBudgets() ([]Budget, error) {
	var budgets []Budget
	result := s.DB.Find(&budgets)
	return budgets, result.Error
}

func (s *Service) AddBudget(name, description, beneficiaryID string, amount Money, intervalMonths int) error {
	return s.DB.Create(&Budget{
		Name:           name,
		Description:    description,
		BeneficiaryID:  beneficiaryID,
		Amount:         amount,
		IntervalMonths: intervalMonths,
	}).Error
}

func (s *Service) UpdateBudget(oldName, newName, description, beneficiaryID string, amount Money, interval int) error {
	if oldName != newName {
		return s.DB.Model(&Budget{Name: oldName}).Updates(Budget{
			Name:           newName,
			Description:    description,
			BeneficiaryID:  beneficiaryID,
			Amount:         amount,
			IntervalMonths: interval,
		}).Error
	}
	return s.DB.Model(&Budget{Name: oldName}).Updates(Budget{
		Description:    description,
		BeneficiaryID:  beneficiaryID,
		Amount:         amount,
		IntervalMonths: interval,
	}).Error
}

func (s *Service) DeleteBudget(name string) error {
	return s.DB.Delete(&Budget{Name: name}).Error
}
