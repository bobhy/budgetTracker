package transactionImport

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"wailts/models"
)

// ParsedTransaction represents a normalized transaction from a CSV file
type ParsedTransaction struct {
	PostedDate  models.Date
	Amount      models.Money
	Description string
	Beneficiary string
	RawHint     string
}

// Parser is the interface that all CSV parsers must implement
type Parser interface {
	Parse(reader io.Reader) ([]ParsedTransaction, error)
}

// GetParser returns the appropriate parser for a given account name
// Using Account Name as the format identifier as per plan
func GetParser(accountName string) (Parser, error) {
	switch accountName {
	case "CapitalOne":
		return &CapitalOneParser{}, nil
	case "WfChecking":
		return &WFCheckingParser{}, nil
	case "WfVisa":
		return &WFVisaParser{}, nil
	default:
		return nil, fmt.Errorf("no parser found for account: %s", accountName)
	}
}

// --- Parsers ---

// CapitalOneParser matches "CapitalOne" format
type CapitalOneParser struct{}

func (p *CapitalOneParser) Parse(reader io.Reader) ([]ParsedTransaction, error) {
	r := csv.NewReader(reader)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Printf("[Parser] CapitalOne CSV read error: %v\n", err)
		return nil, err
	}
	fmt.Printf("[Parser] CapitalOne read %d records (including header)\n", len(records))

	var results []ParsedTransaction

	// Header expected: Transaction Date, Posted Date, Card No., Description, Category, Debit, Credit
	// We skip the first row
	if len(records) > 0 {
		records = records[1:]
	}

	for i, row := range records {
		if len(row) < 7 {
			fmt.Printf("[Parser] Skipping short record at row %d: %v\n", i+1, row)
			continue
		}

		// Posted Date (Col 1)
		postedDate, err := parseDate(row[1])
		if err != nil {
			continue
		}

		// Card No. (Col 2) -> Beneficiary
		cardNo := row[2]
		beneficiary := "Us"
		if strings.HasSuffix(cardNo, "3028") {
			beneficiary = "Bob"
		} else if strings.HasSuffix(cardNo, "6539") {
			beneficiary = "Jessie"
		}

		// Description (Col 3)
		description := row[3]

		// Category (Col 4) -> RawHint
		rawHint := row[4]

		// Debit (Col 5), Credit (Col 6) -> Amount
		// Formula: (Debit*100) - (Credit*100)
		debitVal := parseAmountVal(row[5])
		creditVal := parseAmountVal(row[6])
		amount := models.Money(int64((debitVal * 100) - (creditVal * 100)))

		results = append(results, ParsedTransaction{
			PostedDate:  postedDate,
			Amount:      amount,
			Description: description,
			Beneficiary: beneficiary,
			RawHint:     rawHint,
		})
	}
	return results, nil
}

// WFCheckingParser matches "WfChecking" format
type WFCheckingParser struct{}

func (p *WFCheckingParser) Parse(reader io.Reader) ([]ParsedTransaction, error) {
	return parseWellsFargo(reader)
}

// WFVisaParser matches "WfVisa" format
type WFVisaParser struct{}

func (p *WFVisaParser) Parse(reader io.Reader) ([]ParsedTransaction, error) {
	return parseWellsFargo(reader)
}

// Wrapper for WF logic since Checking and Visa share identical structure and logic in the current rule
func parseWellsFargo(reader io.Reader) ([]ParsedTransaction, error) {
	r := csv.NewReader(reader)
	r.FieldsPerRecord = -1 // Allow variable fields
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var results []ParsedTransaction

	// No Header
	for _, row := range records {
		if len(row) < 5 {
			continue
		}

		// Col 0: Posted Date
		postedDate, err := parseDate(row[0])
		if err != nil {
			continue
		}

		// Col 1: Amount
		// Formula: -1 * value * 100
		val := parseAmountVal(row[1])
		amount := models.Money(int64(-1 * val * 100))

		// Col 2, 3 ignored

		// Col 4: Description
		// Note usage of row[4] matches 5th column 0-indexed?
		// "field 5: maps to .Description" -> index 4.
		// "field 1: maps to .PostedDate" -> index 0.
		description := row[4]

		results = append(results, ParsedTransaction{
			PostedDate:  postedDate,
			Amount:      amount,
			Description: description,
			Beneficiary: "Us",
			RawHint:     "",
		})
	}
	return results, nil
}

// --- Helpers ---

func parseDate(s string) (models.Date, error) {
	layouts := []string{"2006-01-02", "01/02/2006"}
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return models.Date(t.Format("2006-01-02")), nil
		}
	}
	return "", fmt.Errorf("unable to parse date: %s", s)
}

func parseAmountVal(s string) float64 {
	if s == "" {
		return 0
	}
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, " ", "") // Shouldnt happen but safe
	val, _ := strconv.ParseFloat(s, 64)
	return val
}
