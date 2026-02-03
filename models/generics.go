package models

import (
	"strings"

	"gorm.io/gorm"
)

// GetAll returns all records of type T matching the optional conditions.
func GetAll[T any](db *gorm.DB, conds ...interface{}) ([]T, error) {
	var items []T
	err := db.Find(&items, conds...).Error
	return items, err
}

// GetByID returns a record of type T by its ID.
func GetByID[T any](db *gorm.DB, id interface{}) (*T, error) {
	var item T
	if err := db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Create inserts a new record.
func Create[T any](db *gorm.DB, item *T) error {
	return db.Create(item).Error
}

// Update updates a record.
func Update[T any](db *gorm.DB, item *T) error {
	return db.Save(item).Error
}

// Delete deletes a record.
func Delete[T any](db *gorm.DB, item *T) error {
	return db.Delete(item).Error
}

// GetPage returns a paginated list of records of type T.
// offset: number of records to skip
// count: number of records to return
// order: sort order string (e.g. "created_at desc")
// query: optional filter query (e.g. "name LIKE ?")
// args: arguments for the query
func GetPage[T any](db *gorm.DB, skip, count int, order string, query interface{}, args ...interface{}) ([]T, int64, error) {
	var items []T
	var total int64

	dbQuery := db.Model(new(T))

	if query != nil {
		dbQuery = dbQuery.Where(query, args...)
	}

	// Count total records matching filter before pagination
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if order != "" {
		dbQuery = dbQuery.Order(order)
	}

	err := dbQuery.Offset(skip).Limit(count).Find(&items).Error
	return items, total, err
}

// BuildOrderString converts slice of SortOption to a GORM order string.
func BuildOrderString(sortKeys []SortOption) string {
	var parts []string
	for _, sk := range sortKeys {
		if sk.Key == "" || sk.Key == "none" {
			continue
		}

		column := sk.Key

		direction := "asc"
		if sk.Direction == "desc" {
			direction = "desc"
		}
		parts = append(parts, column+" "+direction)
	}
	return strings.Join(parts, ", ")
}
