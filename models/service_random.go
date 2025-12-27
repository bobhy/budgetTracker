package models

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *Service) GenerateBeneficiaries(count int) error {
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("Beneficiary %s", randomString(5))
		if err := s.AddBeneficiary(name); err != nil {
			return err
		}
	}
	return nil
}
