package transactionImport

import (
	"fmt"
	"wailts/models"

	"gorm.io/gorm"
)

// ProcessRaw imports parsed transactions into the RawTransaction table
func ProcessRaw(db *gorm.DB, account string, transactions []ParsedTransaction) error {
	// 1. Fetch existing Transactions for this account to determine "add" vs "update"
	// Optimization: we could filter by date range of the new transactions
	var existingTransactions []models.Transaction
	if err := db.Where("account = ?", account).Find(&existingTransactions).Error; err != nil {
		return err
	}

	// Index existing transactions for fast lookup
	// Key: Date|Amount|Description (Account is fixed)
	transactionMap := make(map[string]bool)
	for _, t := range existingTransactions {
		key := generateKey(t.PostedDate, t.Amount, t.Description)
		transactionMap[key] = true
	}
	fmt.Printf("[Processor] Loaded %d existing transactions for account %s\n", len(existingTransactions), account)

	// 2. Fetch existing RawTransactions to ensure idempotency (update instead of duplicate)
	var existingRaw []models.RawTransaction
	if err := db.Where("account = ?", account).Find(&existingRaw).Error; err != nil {
		return err
	}

	rawMap := make(map[string]uint) // Key -> ID
	for _, r := range existingRaw {
		key := generateKey(r.PostedDate, r.Amount, r.Description)
		rawMap[key] = r.ID
	}
	fmt.Printf("[Processor] Loaded %d existing raw transactions\n", len(existingRaw))

	// 3. Process each parsed transaction
	added := 0
	updated := 0
	for _, pt := range transactions {
		key := generateKey(pt.PostedDate, pt.Amount, pt.Description)

		// Determine Action
		action := "add"
		if transactionMap[key] {
			action = "update"
		}

		// Prepare model
		raw := models.RawTransaction{
			PostedDate:  pt.PostedDate,
			Account:     account,
			Amount:      pt.Amount,
			Description: pt.Description,
			Beneficiary: pt.Beneficiary,
			RawHint:     pt.RawHint,
			Action:      action,
			// Budget: is imported as empty string, user must set to somethingh non-empty to load into transactions.
		}

		// Check if we already have this in Raw
		if rawID, exists := rawMap[key]; exists {
			// Update existing Raw record
			raw.ID = rawID
			if err := db.Save(&raw).Error; err != nil {
				return err
			}
			updated++
		} else {
			// Create new
			if err := db.Create(&raw).Error; err != nil {
				return err
			}
			rawMap[key] = raw.ID
			added++
		}
	}
	fmt.Printf("[Processor] Processing complete. Added: %d, Updated: %d\n", added, updated)

	return nil
}

func generateKey(date models.Date, amount models.Money, desc string) string {
	return fmt.Sprintf("%s|%d|%s", date, amount, desc)
}
