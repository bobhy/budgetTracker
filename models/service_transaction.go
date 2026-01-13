package models

func (s *Service) GetTransactions() ([]Transaction, error) {
	var transactions []Transaction
	// Preload Account?
	result := s.DB.Preload("Account").Find(&transactions)
	return transactions, result.Error
}

func (s *Service) GetTransactionsPaginated(start, count int, sortKeys []SortOption) ([]Transaction, error) {
	var transactions []Transaction
	query := s.DB.Preload("Account").Model(&Transaction{})

	// Sorting
	if len(sortKeys) > 0 {
		for _, sk := range sortKeys {
			if sk.Key != "" && sk.Key != "none" {
				// Convert frontend key to DB column
				// e.g. "PostedDate" -> "posted_date"
				// Gorm handles CamelCase to snake_case usually if using struct fields?
				// But we probably get "PostedDate" from frontend.
				// Safe mapping or direct usage?
				// Gorm Order("posted_date desc")
				// We need to map or trust.
				// Let's rely on GORM's NamingStrategy or manual mapping.
				// Simple snake_case conversion or switch:
				column := ""
				switch sk.Key {
				case "ID":
					column = "id"
				case "PostedDate":
					column = "posted_date"
				case "AccountID":
					column = "account_id"
				case "Amount":
					column = "amount"
				case "Description":
					column = "description"
				case "Beneficiary":
					column = "beneficiary"
				case "BudgetLine":
					column = "budget_line"
				case "Tag":
					column = "tag"
				}

				if column != "" {
					direction := "asc"
					if sk.Direction == "desc" {
						direction = "desc"
					}
					query = query.Order(column + " " + direction)
				}
			}
		}
	} else {
		// Default sort
		//query = query.Order("posted_date desc")
	}

	// Pagination
	err := query.Offset(start).Limit(count).Find(&transactions).Error
	return transactions, err
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
