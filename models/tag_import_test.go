package models

import (
	"testing"
	"time"
)

func SetupTestService(t *testing.T) *Service {
	// Mode=memory (shared cache matches generics_test.go)
	// We pass "file::memory:" so NewService appends "?_foreign_keys=on" -> "file::memory:?_foreign_keys=on"
	s, err := NewService("file::memory:")
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	return s
}

func TestApplyTags(t *testing.T) {
	s := SetupTestService(t)

	// 1. Setup Data
	// - Tag: "groceries" -> Budget: "Food"
	// - RawTransaction: "purchase authorized on 01/01 kroger" -> Expected Tag: "kroger" -> No Budget
	// - RawTransaction: "purchase authorized on 01/01 whole foods" -> Expected Tag: "whole foods" -> Budget: "Food" (needs tag entry)

	// Note: seed() in NewService adds placeholder data. We can just add to it.

	// Create Beneficiary
	ben := Beneficiary{Name: "Us"}
	if err := s.AddBeneficiary(&ben); err != nil {
		t.Fatalf("Failed to add beneficiary: %v", err)
	}

	// Create Budgets
	foodBudget := Budget{Name: "Food", Beneficiary: "Us"}
	if err := s.AddBudget(&foodBudget); err != nil {
		t.Fatalf("Failed to add budget: %v", err)
	}

	// Create Tags
	// Map "whole foods" to "Food" budget
	tag := Tag{Name: "whole foods", Budget: "Food"}
	if err := s.AddTag(&tag); err != nil {
		t.Fatalf("Failed to add tag: %v", err)
	}

	// Create Raw Transactions
	now := Date(time.Now().Format("2006-01-02"))
	rawTxs := []RawTransaction{
		{
			Description: "purchase authorized on 01/12 whole foods store #123",
			Amount:      10000,
			Account:     "TestAcc",
			PostedDate:  now,
		},
		{
			Description: "purchase authorized on 01/12 unknown store",
			Amount:      5000,
			Account:     "TestAcc",
			PostedDate:  now,
		},
		{
			Description: "money transfer authorized on 01/12 zelle to friend",
			Amount:      2000,
			Account:     "TestAcc",
			PostedDate:  now,
		},
	}

	for _, rt := range rawTxs {
		if err := s.AddRawTransaction(&rt); err != nil {
			t.Fatalf("Failed to add raw tx: %v", err)
		}
	}

	// 2. Run ApplyTags
	count, err := s.ApplyTags()
	if err != nil {
		t.Fatalf("ApplyTags failed: %v", err)
	}

	// 3. Verify Results

	// Check count: should be 1 (whole foods -> Food)
	// Actually, update count depends on how many rows were updated by the second query (Budget Mapping).
	// The "whole foods" one should match.
	if count != 1 {
		t.Errorf("Expected 1 updated row, got %d", count)
	}

	updatedTxs, err := s.GetRawTransactions()
	if err != nil {
		t.Fatalf("Failed to get raw txs: %v", err)
	}

	for _, tx := range updatedTxs {
		t.Logf("Tx: %s, Tag: '%s', Budget: '%s'", tx.Description, tx.Tag, tx.Budget)

		if tx.Description == "purchase authorized on 01/12 whole foods store #123" {
			if tx.Tag != "whole foods store #123" {
				t.Errorf("Expected tag 'whole foods store #123', got '%s'", tx.Tag)
			}
			if tx.Budget != "Food" {
				t.Errorf("Expected budget 'Food', got '%s'", tx.Budget)
			}
		} else if tx.Description == "purchase authorized on 01/12 unknown store" {
			if tx.Tag != "unknown store" {
				t.Errorf("Expected tag 'unknown store', got '%s'", tx.Tag)
			}
			if tx.Budget != "" {
				t.Errorf("Expected empty budget, got '%s'", tx.Budget)
			}
		} else if tx.Description == "money transfer authorized on 01/12 zelle to friend" {
			if tx.Tag != "friend" {
				t.Errorf("Expected tag 'friend', got '%s'", tx.Tag)
			}
		}
	}
}
