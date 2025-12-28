package models

func (s *Service) GetTransactions() ([]Transaction, error) {
	var transactions []Transaction
	// Preload Account?
	result := s.DB.Preload("Account").Find(&transactions)
	return transactions, result.Error
}

func (s *Service) AddTransaction(postedDate Date, accountID string, amount Money, description, tag, beneficiary, budgetLine, rawHint string) error {
	return s.DB.Create(&Transaction{
		PostedDate:  postedDate,
		AccountID:   accountID,
		Amount:      amount,
		Description: description,
		Tag:         tag,
		Beneficiary: beneficiary,
		BudgetLine:  budgetLine,
		RawHint:     rawHint,
	}).Error
}

// UpdateTransaction uses ID for identification
func (s *Service) UpdateTransaction(id uint, postedDate Date, accountID string, amount Money, description, tag, beneficiary, budgetLine, rawHint string) error {
	return s.DB.Model(&Transaction{}).Where("id = ?", id).Updates(map[string]interface{}{
		"posted_date": postedDate,
		"account_id":  accountID,
		"amount":      amount,
		"description": description,
		"tag":         tag,
		"beneficiary": beneficiary,
		"budget_line": budgetLine,
		"raw_hint":    rawHint,
	}).Error
}

func (s *Service) DeleteTransaction(id uint) error {
	return s.DB.Delete(&Transaction{}, id).Error
}
