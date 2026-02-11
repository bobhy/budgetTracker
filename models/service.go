package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

var allTables = []any{
	&Beneficiary{},
	&Account{},
	&Budget{},
	&Tag{},
	&Transaction{},
	&RawTransaction{},
}

func NewService(dbPath string) (*Service, error) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s?%s", dbPath, "_foreign_keys=on")), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migrate
	err = db.AutoMigrate(allTables...)
	if err != nil {
		return nil, err
	}

	return &Service{DB: db}, nil
}

// Clean drops all tables and re-migrates them, then seeds production data
func (s *Service) Clean() error {
	// Drop tables
	err := s.DB.Migrator().DropTable(allTables...)
	if err != nil {
		return err
	}
	// Re-migrate
	err = s.DB.AutoMigrate(allTables...)
	if err != nil {
		return err
	}

	return s.seed()
}

func seedTable[T any](s *Service, records []T) error {
	for _, r := range records {
		if err := Create(s.DB, &r); err != nil {
			return err
		}
	}
	return nil
}

/// Placeholder user must update before a RawTransaction can be moved to Transactions table.

const UNCATEGORIZED_BUDGET = ""             // user hasn't assigned any budget -- tx stays in raw table
const PLACEHOLDER_BUDGET = "--unbudgeted--" // a budget the user can assign to complete the import but worry about budget details later
const PLACEHOLDER_BENEFICIARY = "--none--"

func (s *Service) seed() error {

	// Seed Beneficiaries
	err := seedTable(s, []Beneficiary{
		{Name: PLACEHOLDER_BENEFICIARY},
		{Name: "Us"},
		{Name: "Bob"},
		{Name: "Jessie"},
	})

	if err != nil {
		return fmt.Errorf("failed to seed beneficiaries: %w", err)
	}

	// Seed Accounts
	err = seedTable(s, []Account{
		{
			Name:        "CapitalOne",
			Description: "Capital One rewards Credit Account",
			Beneficiary: "Us",
		},
		{
			Name:        "WfChecking",
			Description: "Wells Fargo checking",
			Beneficiary: "Us",
		},
		{
			Name:        "WfVisa",
			Description: "Wells Fargo Visa",
			Beneficiary: "Us",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to seed accounts: %w", err)
	}

	// Seed Budgets
	err = seedTable(s, []Budget{
		{
			Name:        PLACEHOLDER_BUDGET,
			Description: "Update to move raw transaction into transactions",
			Beneficiary: PLACEHOLDER_BENEFICIARY,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to seed budgets: %w", err)
	}

	// Seed Tags
	err = seedTable(s, []Tag{
		{
			Name:   "",
			Budget: PLACEHOLDER_BUDGET,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to seed tags: %w", err)
	}

	return nil
}

// Templateized CRUD functions for all the tables
// todo: see whether typescript can bind to underlying generic db calls?

// --- Beneficiaries ---

func (s *Service) GetBeneficiaries() ([]Beneficiary, error) {
	return GetAll[Beneficiary](s.DB)
}

func (s *Service) GetBeneficiariesPaginated(start, count int, sortKeys []SortOption) ([]Beneficiary, error) {
	orderStr := BuildOrderString(sortKeys)
	beneficiaries, _, err := GetPage[Beneficiary](s.DB, start, count, orderStr, nil)
	return beneficiaries, err
}

func (s *Service) AddBeneficiary(beneficiary *Beneficiary) error {
	return Create(s.DB, beneficiary)
}

func (s *Service) UpdateBeneficiary(oldBeneficiary, newBeneficiary *Beneficiary) error {
	return s.DB.Model(oldBeneficiary).Updates(newBeneficiary).Error
}

func (s *Service) DeleteBeneficiary(beneficiary *Beneficiary) error {
	return Delete(s.DB, beneficiary)
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

func (s *Service) AddAccount(account *Account) error {
	return Create(s.DB, account)
}

func (s *Service) UpdateAccount(oldAccount, newAccount *Account) error {
	return s.DB.Model(oldAccount).Updates(newAccount).Error
}

func (s *Service) DeleteAccount(account *Account) error {
	return Delete(s.DB, account)
}

// --- Budgets ---

func (s *Service) GetBudgets() ([]Budget, error) {
	return GetAll[Budget](s.DB)
}

func (s *Service) GetBudgetsPaginated(start, count int, sortKeys []SortOption) ([]Budget, error) {
	orderStr := BuildOrderString(sortKeys)
	budgets, _, err := GetPage[Budget](s.DB, start, count, orderStr, nil)
	return budgets, err
}

func (s *Service) AddBudget(budget *Budget) error {
	return Create(s.DB, budget)
}

func (s *Service) UpdateBudget(oldBudget, newBudget *Budget) error {
	return s.DB.Model(oldBudget).Updates(newBudget).Error
}

func (s *Service) DeleteBudget(budget *Budget) error {
	return Delete(s.DB, budget)
}

// --- Tags ---

func (s *Service) GetTags() ([]Tag, error) {
	return GetAll[Tag](s.DB)
}

func (s *Service) GetTagsPaginated(start, count int, sortKeys []SortOption) ([]Tag, error) {
	orderStr := BuildOrderString(sortKeys)
	Tags, _, err := GetPage[Tag](s.DB, start, count, orderStr, nil)
	return Tags, err
}

func (s *Service) AddTag(Tag *Tag) error {
	return Create(s.DB, Tag)
}

func (s *Service) UpdateTag(oldTag, newTag *Tag) error {
	return s.DB.Model(oldTag).Updates(newTag).Error
}

func (s *Service) DeleteTag(Tag *Tag) error {
	return Delete(s.DB, Tag)
}

// --- Transactions ---

func (s *Service) GetTransactions() ([]Transaction, error) {
	return GetAll[Transaction](s.DB)
}

func (s *Service) GetTransactionsPaginated(start, count int, sortKeys []SortOption) ([]Transaction, error) {
	orderStr := BuildOrderString(sortKeys)
	txs, _, err := GetPage[Transaction](s.DB, start, count, orderStr, nil)
	return txs, err
}

func (s *Service) AddTransaction(transaction *Transaction) error {
	return Create(s.DB, transaction)
}

func (s *Service) UpdateTransaction(oldTransaction, newTransaction *Transaction) error {
	return s.DB.Model(oldTransaction).Updates(newTransaction).Error
}

func (s *Service) DeleteTransaction(transaction *Transaction) error {
	return Delete(s.DB, transaction)
}

// --- Raw Transactions ---

func (s *Service) GetRawTransactions() ([]RawTransaction, error) {
	return GetAll[RawTransaction](s.DB)
}
func (s *Service) GetRawTransactionsPaginated(start, count int, sortKeys []SortOption) ([]RawTransaction, error) {
	orderStr := BuildOrderString(sortKeys)
	rawTxs, _, err := GetPage[RawTransaction](s.DB, start, count, orderStr, nil)
	return rawTxs, err
}

func (s *Service) GetRawTransactionCount() (int64, error) {
	var count int64
	err := s.DB.Model(&RawTransaction{}).Count(&count).Error
	return count, err
}

func (s *Service) AddRawTransaction(rawTransaction *RawTransaction) error {
	return Create(s.DB, rawTransaction)
}

func (s *Service) UpdateRawTransaction(oldRawTransaction, newRawTransaction *RawTransaction) error {
	return s.DB.Model(oldRawTransaction).Updates(newRawTransaction).Error
}

func (s *Service) DeleteRawTransaction(rawTransaction *RawTransaction) error {
	return Delete(s.DB, rawTransaction)
}

func (s *Service) FinalizeImport() (string, error) {
	var rawList []RawTransaction
	if err := s.DB.Find(&rawList).Error; err != nil {
		return "", err
	}

	added := 0   //new tx in tx table
	updated := 0 // existing tx updated in tx table
	skipped := 0 // raw tx left in raw table because uncategorized.

	tx := s.DB.Begin()

	for _, raw := range rawList {
		if raw.Budget == UNCATEGORIZED_BUDGET {
			skipped++
			continue
		}

		switch raw.Action {
		case "add":
			// Create new Transaction
			t := Transaction{
				PostedDate:  raw.PostedDate,
				Account:     raw.Account,
				Amount:      raw.Amount,
				Description: raw.Description,
				Beneficiary: raw.Beneficiary,
				Budget:      raw.Budget,
				RawHint:     raw.RawHint,
			}
			if err := tx.Create(&t).Error; err != nil {
				tx.Rollback()
				return "", err
			}
			added++
		case "update":
			var target Transaction
			result := tx.Where("account_id = ? AND posted_date = ? AND amount = ? AND description = ?",
				raw.Account, raw.PostedDate, raw.Amount, raw.Description).First(&target)

			if result.Error == nil {
				// Found match. Update it.
				target.Beneficiary = raw.Beneficiary
				target.Budget = raw.Budget
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
					Account:     raw.Account,
					Amount:      raw.Amount,
					Description: raw.Description,
					Beneficiary: raw.Beneficiary,
					Budget:      raw.Budget,
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
	if err := tx.Exec("DELETE FROM raw_transactions WHERE budget != ?", UNCATEGORIZED_BUDGET).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()
	return fmt.Sprintf("Finalized: %d added, %d updated, %d remain to be categorized.", added, updated, skipped), nil
}

func (s *Service) ApplyTags() (int64, error) {
	// 1. Tagging Query
	// Updates raw_transactions.tag based on description patterns
	taggingQuery := `
	WITH t1 AS (
		SELECT
			id,
			CASE
				WHEN description LIKE 'money transfer authorized on __/__ %' THEN substr(description, 36)
				WHEN description LIKE 'purchase authorized on __/__ %' THEN substr(description, 30)
				WHEN description LIKE 'purchase intl authorized on __/__ %' THEN substr(description, 35)
				ELSE description
			END AS stem1
		FROM raw_transactions
	),
	t2 AS (
		SELECT
			id,
			CASE
				WHEN stem1 LIKE '___*%' THEN substr(stem1, 5)
				WHEN stem1 LIKE 'cash app*%' THEN substr(stem1, 10)
				WHEN stem1 LIKE 'zelle to %' THEN substr(stem1, 10)
				WHEN stem1 LIKE 'paypal *%' THEN substr(stem1, 9)
				ELSE stem1
			END AS stem
		FROM t1
	)
	UPDATE raw_transactions
	SET tag = t2.stem
	FROM t2
	WHERE raw_transactions.id = t2.id;
	`

	// 2. Budget Mapping Query
	// Updates raw_transactions.budget based on tag matches in tags table
	budgetQuery := `
	UPDATE raw_transactions
	SET budget = t.budget
	FROM tags t
	WHERE substr(raw_transactions.tag, 1, length(t.name)) = t.name;
	`

	tx := s.DB.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	// Run Tagging
	if err := tx.Exec(taggingQuery).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("tagging query failed: %w", err)
	}

	// Run Budget Mapping
	result := tx.Exec(budgetQuery)
	if result.Error != nil {
		tx.Rollback()
		return 0, fmt.Errorf("budget mapping query failed: %w", result.Error)
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return result.RowsAffected, nil
}
