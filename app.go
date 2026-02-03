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
