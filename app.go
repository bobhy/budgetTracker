package main

import (
	"context"
	"wailts/models"
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
func (a *App) GenerateBeneficiaries(count int) error {
	return a.service.GenerateBeneficiaries(count)
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
func (a *App) GenerateAccounts(count int) error {
	return a.service.GenerateAccounts(count)
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
func (a *App) GenerateBudgets(count int) error {
	return a.service.GenerateBudgets(count)
}

// --- Transactions ---

func (a *App) GetTransactions() ([]models.Transaction, error) {
	return a.service.GetTransactions()
}
func (a *App) AddTransaction(postedDate models.Date, accountID string, amount models.Money, description, tag string) error {
	return a.service.AddTransaction(postedDate, accountID, amount, description, tag)
}
func (a *App) UpdateTransaction(id uint, postedDate models.Date, accountID string, amount models.Money, description, tag string) error {
	return a.service.UpdateTransaction(id, postedDate, accountID, amount, description, tag)
}
func (a *App) DeleteTransaction(id uint) error {
	return a.service.DeleteTransaction(id)
}
func (a *App) GenerateTransactions(count int) error {
	return a.service.GenerateTransactions(count)
}
