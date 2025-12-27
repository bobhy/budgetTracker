package models

import (
	"fmt"
	"math/rand"
)

func (s *Service) GetBudgets() ([]Budget, error) {
	var budgets []Budget
	result := s.db.Find(&budgets)
	return budgets, result.Error
}

func (s *Service) AddBudget(name, description, beneficiaryID string, amount Money, intervalMonths int) error {
	return s.db.Create(&Budget{
		Name:           name,
		Description:    description,
		BeneficiaryID:  beneficiaryID,
		Amount:         amount,
		IntervalMonths: intervalMonths,
	}).Error
}

func (s *Service) UpdateBudget(oldName, newName, description, beneficiaryID string, amount Money, interval int) error {
	if oldName != newName {
		return s.db.Model(&Budget{Name: oldName}).Updates(Budget{
			Name:           newName,
			Description:    description,
			BeneficiaryID:  beneficiaryID,
			Amount:         amount,
			IntervalMonths: interval,
		}).Error
	}
	return s.db.Model(&Budget{Name: oldName}).Updates(Budget{
		Description:    description,
		BeneficiaryID:  beneficiaryID,
		Amount:         amount,
		IntervalMonths: interval,
	}).Error
}

func (s *Service) DeleteBudget(name string) error {
	return s.db.Delete(&Budget{Name: name}).Error
}

func (s *Service) GenerateBudgets(count int) error {
	bens, err := s.GetBeneficiaries()
	if err != nil {
		return err
	}
	if len(bens) == 0 {
		return fmt.Errorf("no beneficiaries found")
	}

	for i := 0; i < count; i++ {
		ben := bens[rand.Intn(len(bens))]
		name := fmt.Sprintf("Budget %s", randomString(5))
		amount := Money(rand.Int63n(100000)) // 0 to 1000.00
		interval := rand.Intn(12) + 1        // 1 to 12 months
		if err := s.AddBudget(name, "Random budget", ben.Name, amount, interval); err != nil {
			return err
		}
	}
	return nil
}
