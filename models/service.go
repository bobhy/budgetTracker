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

// --- Beneficiaries ---

func (s *Service) GetBeneficiaries() ([]Beneficiary, error) {
	return GetAll[Beneficiary](s.DB)
}

func (s *Service) AddBeneficiary(name string) error {
	return Create(s.DB, &Beneficiary{Name: name})
}

func (s *Service) UpdateBeneficiary(oldName, newName string) error {
	if oldName == newName {
		return nil
	}
	return s.DB.Model(&Beneficiary{Name: oldName}).Update("name", newName).Error
}

func (s *Service) DeleteBeneficiary(name string) error {
	return Delete(s.DB, &Beneficiary{Name: name})
}

// --- Accounts ---

func (s *Service) GetAccounts() ([]Account, error) {
	return GetAll[Account](s.DB)
}

func (s *Service) GetAccountsPaginated(start, count int, sortKeys []SortOption) ([]Account, error) {
	orderStr := BuildOrderString(sortKeys)
	accounts, _, err := GetPage[Account](s.DB, start, count, orderStr, nil)
	return accounts, err
}

func (s *Service) AddAccount(name, description, beneficiaryID string) error {
	return Create(s.DB, &Account{
		Name:          name,
		Description:   description,
		BeneficiaryID: beneficiaryID,
	})
}

func (s *Service) UpdateAccount(oldName, newName, description, beneficiaryID string) error {
	if oldName != newName {
		return s.DB.Model(&Account{Name: oldName}).Updates(Account{
			Name:          newName,
			Description:   description,
			BeneficiaryID: beneficiaryID,
		}).Error
	}
	return s.DB.Model(&Account{Name: oldName}).Updates(Account{
		Description:   description,
		BeneficiaryID: beneficiaryID,
	}).Error
}

func (s *Service) DeleteAccount(name string) error {
	return Delete(s.DB, &Account{Name: name})
}

// --- Budgets ---

func (s *Service) GetBudgets() ([]Budget, error) {
	return GetAll[Budget](s.DB)
}

func (s *Service) AddBudget(name, description, beneficiaryID string, amount Money, intervalMonths int) error {
	return Create(s.DB, &Budget{
		Name:           name,
		Description:    description,
		BeneficiaryID:  beneficiaryID,
		Amount:         amount,
		IntervalMonths: intervalMonths,
	})
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
	return Delete(s.DB, &Budget{Name: name})
}

// --- Transactions ---

func (s *Service) GetTransactions() ([]Transaction, error) {
	return GetAll[Transaction](s.DB.Preload("Account"))
}

func (s *Service) GetTransactionsPaginated(start, count int, sortKeys []SortOption) ([]Transaction, error) {
	orderStr := BuildOrderString(sortKeys)
	txs, _, err := GetPage[Transaction](s.DB.Preload("Account"), start, count, orderStr, nil)
	return txs, err
}

func (s *Service) AddTransaction(postedDate Date, accountID string, amount Money, description, tag, beneficiary, budgetLine, rawHint string) error {
	return Create(s.DB, &Transaction{
		PostedDate:  postedDate,
		AccountID:   accountID,
		Amount:      amount,
		Description: description,
		Tag:         tag,
		Beneficiary: beneficiary,
		BudgetLine:  budgetLine,
		RawHint:     rawHint,
	})
}

func (s *Service) UpdateTransaction(id uint, postedDate Date, accountID string, amount Money, description, tag, beneficiary, budgetLine, rawHint string) error {
	return s.DB.Model(&Transaction{ID: id}).Updates(map[string]interface{}{
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
	return Delete(s.DB, &Transaction{ID: id})
}

// --- Raw Transactions ---

func (s *Service) GetRawTransactions() ([]RawTransaction, error) {
	var raw []RawTransaction
	err := s.DB.Order("posted_date desc").Find(&raw).Error
	return raw, err
}

func (s *Service) UpdateRawTransaction(id uint, postedDate Date, amount Money, description, beneficiary, budgetLine, rawHint string) error {
	var raw RawTransaction
	if err := s.DB.First(&raw, id).Error; err != nil {
		return err
	}

	raw.PostedDate = postedDate
	raw.Amount = amount
	raw.Description = description
	raw.Beneficiary = beneficiary
	raw.BudgetLine = budgetLine
	raw.RawHint = rawHint

	return s.DB.Save(&raw).Error
}

func (s *Service) DeleteRawTransaction(id uint) error {
	return s.DB.Delete(&RawTransaction{}, id).Error
}

func (s *Service) FinalizeImport() (string, error) {
	var rawList []RawTransaction
	if err := s.DB.Find(&rawList).Error; err != nil {
		return "", err
	}

	added := 0
	updated := 0

	tx := s.DB.Begin()

	for _, raw := range rawList {
		if raw.Action == "add" {
			// Create new Transaction
			t := Transaction{
				PostedDate:  raw.PostedDate,
				AccountID:   raw.AccountID,
				Amount:      raw.Amount,
				Description: raw.Description,
				Beneficiary: raw.Beneficiary,
				BudgetLine:  raw.BudgetLine,
				RawHint:     raw.RawHint,
			}
			if err := tx.Create(&t).Error; err != nil {
				tx.Rollback()
				return "", err
			}
			added++
		} else if raw.Action == "update" {
			var target Transaction
			result := tx.Where("account_id = ? AND posted_date = ? AND amount = ? AND description = ?",
				raw.AccountID, raw.PostedDate, raw.Amount, raw.Description).First(&target)

			if result.Error == nil {
				// Found match. Update it.
				target.Beneficiary = raw.Beneficiary
				target.BudgetLine = raw.BudgetLine
				target.RawHint = raw.RawHint
				if err := tx.Save(&target).Error; err != nil {
					tx.Rollback()
					return "", err
				}
				updated++
			} else {
				// Not found. Treat as new to avoid data loss.
				t := Transaction{
					PostedDate:  raw.PostedDate,
					AccountID:   raw.AccountID,
					Amount:      raw.Amount,
					Description: raw.Description,
					Beneficiary: raw.Beneficiary,
					BudgetLine:  raw.BudgetLine,
					RawHint:     raw.RawHint,
				}
				if err := tx.Create(&t).Error; err != nil {
					tx.Rollback()
					return "", err
				}
				updated++
			}
		}
	}

	// Empty Raw
	if err := tx.Exec("DELETE FROM raw_transactions").Error; err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()
	return fmt.Sprintf("Finalized: %d added, %d updated", added, updated), nil
}
