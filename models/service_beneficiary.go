package models

// GetBeneficiaries returns all beneficiaries
func (s *Service) GetBeneficiaries() ([]Beneficiary, error) {
	var beneficiaries []Beneficiary
	result := s.DB.Find(&beneficiaries)
	return beneficiaries, result.Error
}

// AddBeneficiary adds a new beneficiary
func (s *Service) AddBeneficiary(name string) error {
	return s.DB.Create(&Beneficiary{Name: name}).Error
}

// UpdateBeneficiary updates a beneficiary's name.
// Since name is the primary key, we might need to handle this carefully if it was referenced,
// but for v0.1 we rely on GORM/SQLite behavior or simple updates.
// Note: Updating a primary key is usually 'Model(old).Update("name", new)'
func (s *Service) UpdateBeneficiary(oldName, newName string) error {
	// If name hasn't changed, do nothing
	if oldName == newName {
		return nil
	}
	// For SQLite, updating a PK usually cascades if FK constraints are set with ON UPDATE CASCADE.
	// GORM doesn't default to ON UPDATE CASCADE for all drivers easily without tags.
	// But let's try the standard update.
	return s.DB.Model(&Beneficiary{Name: oldName}).Update("name", newName).Error
}

// DeleteBeneficiary deletes a beneficiary
func (s *Service) DeleteBeneficiary(name string) error {
	return s.DB.Delete(&Beneficiary{Name: name}).Error
}
