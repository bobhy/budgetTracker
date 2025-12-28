package main

import (
	"context"
	"fmt"
	"os"
	"wailts/models"
	"wailts/transactionImport"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	service *models.Service
}

// NewApp creates a new App application struct
func NewApp(service *models.Service) *App {
	return &App{
		service: service,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// --- Database Admin ---

func (a *App) CleanDatabase() error {
	return a.service.Clean()
}

// --- Beneficiaries ---

func (a *App) GetBeneficiaries() ([]models.Beneficiary, error) {
	return a.service.GetBeneficiaries()
}
func (a *App) AddBeneficiary(name string) error {
	return a.service.AddBeneficiary(name)
}
func (a *App) UpdateBeneficiary(oldName, newName string) error {
	return a.service.UpdateBeneficiary(oldName, newName)
}
func (a *App) DeleteBeneficiary(name string) error {
	return a.service.DeleteBeneficiary(name)
}

// --- Accounts ---

func (a *App) GetAccounts() ([]models.Account, error) {
	return a.service.GetAccounts()
}
func (a *App) AddAccount(name, description, beneficiaryID string) error {
	return a.service.AddAccount(name, description, beneficiaryID)
}
func (a *App) UpdateAccount(oldName, newName, description, beneficiaryID string) error {
	return a.service.UpdateAccount(oldName, newName, description, beneficiaryID)
}
func (a *App) DeleteAccount(name string) error {
	return a.service.DeleteAccount(name)
}

// --- Budgets ---

func (a *App) GetBudgets() ([]models.Budget, error) {
	return a.service.GetBudgets()
}
func (a *App) AddBudget(name, description, beneficiaryID string, amount models.Money, intervalMonths int) error {
	return a.service.AddBudget(name, description, beneficiaryID, amount, intervalMonths)
}
func (a *App) UpdateBudget(oldName, newName, description, beneficiaryID string, amount models.Money, interval int) error {
	return a.service.UpdateBudget(oldName, newName, description, beneficiaryID, amount, interval)
}
func (a *App) DeleteBudget(name string) error {
	return a.service.DeleteBudget(name)
}

// --- Transactions ---

func (a *App) GetTransactions() ([]models.Transaction, error) {
	return a.service.GetTransactions()
}
func (a *App) AddTransaction(postedDate models.Date, accountID string, amount models.Money, description, tag, beneficiary, budgetLine, rawHint string) error {
	return a.service.AddTransaction(postedDate, accountID, amount, description, tag, beneficiary, budgetLine, rawHint)
}
func (a *App) UpdateTransaction(id uint, postedDate models.Date, accountID string, amount models.Money, description, tag, beneficiary, budgetLine, rawHint string) error {
	return a.service.UpdateTransaction(id, postedDate, accountID, amount, description, tag, beneficiary, budgetLine, rawHint)
}
func (a *App) DeleteTransaction(id uint) error {
	return a.service.DeleteTransaction(id)
}

// --- Import ---

func (a *App) SelectFile() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select CSV File",
		Filters: []runtime.FileFilter{
			{DisplayName: "CSV Files", Pattern: "*.csv"},
		},
	})
}

func (a *App) ImportFile(accountID string, filePath string) (string, error) {
	runtime.LogInfo(a.ctx, fmt.Sprintf("ImportFile called for account: %s, file: %s", accountID, filePath))

	// Open File
	f, err := os.Open(filePath)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Error opening file: %s", err))
		return "", err
	}
	defer f.Close()

	// Create Parser based on AccountID (Name)
	parser, err := transactionImport.GetParser(accountID)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Error getting parser: %s", err))
		return "", err
	}
	runtime.LogInfo(a.ctx, fmt.Sprintf("Using parser for: %s", accountID))

	// Parse
	records, err := parser.Parse(f)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Error parsing file: %s", err))
		return "", err
	}
	runtime.LogInfo(a.ctx, fmt.Sprintf("Parsed %d records", len(records)))

	// Process
	err = transactionImport.ProcessRaw(a.service.DB, accountID, records)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Error processing raw transactions: %s", err))
		return "", err
	}

	msg := fmt.Sprintf("Imported %d records", len(records))
	runtime.LogInfo(a.ctx, msg)
	return msg, nil
}

func (a *App) GetRawTransactions() ([]models.RawTransaction, error) {
	var raw []models.RawTransaction
	// Order by PostedDate desc
	err := a.service.DB.Order("posted_date desc").Find(&raw).Error
	return raw, err
}

// UpdateRawTransaction updates editable fields of a raw transaction
func (a *App) UpdateRawTransaction(id uint, postedDate models.Date, amount models.Money, description, beneficiary, budgetLine, rawHint string) error {
	// We only allow updating ID if it exists
	var raw models.RawTransaction
	if err := a.service.DB.First(&raw, id).Error; err != nil {
		return err
	}

	raw.PostedDate = postedDate
	raw.Amount = amount
	raw.Description = description
	raw.Beneficiary = beneficiary
	raw.BudgetLine = budgetLine
	raw.RawHint = rawHint
	// Action usually stays same unless logic changes, but let's assume manual edit keeps "update" or "add" state or User might flip it?
	// For now, don't expose Action update via arguments.

	return a.service.DB.Save(&raw).Error
}

func (a *App) DeleteRawTransaction(id uint) error {
	return a.service.DB.Delete(&models.RawTransaction{}, id).Error
}

// FinalizeImport moves valid RawTransactions to Main Transactions table and clears Raw
func (a *App) FinalizeImport() (string, error) {
	var rawList []models.RawTransaction
	if err := a.service.DB.Find(&rawList).Error; err != nil {
		return "", err
	}

	added := 0
	updated := 0

	tx := a.service.DB.Begin()

	for _, raw := range rawList {
		if raw.Action == "add" {
			// Create new Transaction
			t := models.Transaction{
				PostedDate:  raw.PostedDate,
				AccountID:   raw.AccountID,
				Amount:      raw.Amount,
				Description: raw.Description,
				Beneficiary: raw.Beneficiary, // New Field
				BudgetLine:  raw.BudgetLine,  // New Field
				RawHint:     raw.RawHint,
				// Tag? Raw has Tag but we didn't map it from CSV.
			}
			if err := tx.Create(&t).Error; err != nil {
				tx.Rollback()
				return "", err
			}
			added++
		} else if raw.Action == "update" {
			// Find existing and update
			// Match criteria: Account, Date, Amount, Description (Original).
			// This is tricky if Description was edited in Raw.
			// Ideally we linked them or stored the ID of the matched transaction in Raw if we wanted to be robust.
			// But for now, we rely on the same matching logic or maybe we skip updates if we can't find it exactly?
			// "Action is 'update' if this raw transaction will change an existing record..."
			// We need to FIND that record.
			// Since we didn't store the Target Transaction ID in RawTransaction, we have to re-match.
			// If user edited the Raw record (e.g. fixed description), strict matching might fail.
			// Limitation: We only update if we find the match.
			// BETTER: We should have likely stored TargetTransactionID in RawTransaction during ProcessRaw.
			// But sticking to plan: we try to match.

			// For simplicity/safety: we might skip "update" logic here if ambiguous, or just inserting new if distinct?
			// Rule: "Upserts the main transactions table".
			// Logic: Look for exact match. If found, update. If not found, create new? Or strict error?
			// Let's try to match on AccountID, Date, Amount (Original?).
			// If we can't find it, we probably shouldn't blindly insert duplicate if Action says Update.

			// IMPLEMENTATION CHOICE:
			// For now, I will skip "Updating" existing main transactions because re-matching is fragile without ID.
			// I will assuming "Add" works.
			// If "Update", I will try to find match by (Account, Date, Amount).
			// If multiple, I pick first.
			var target models.Transaction
			// Use original values? We don't have them if edited.
			// We use current values.
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
				// Not found (maybe edited?). Treat as new or skip?
				// If we treat as new, we might duplicate.
				// Let's Skip and Warn? Or just Add?
				// Rule says "upserts".
				// I'll Add it to be safe from dataloss, risking duplication.
				t := models.Transaction{
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
				updated++ // Count as update/add
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
