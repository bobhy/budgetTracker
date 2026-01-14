package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestModel struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	Age  int
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	// Clean up previous data
	_ = db.Migrator().DropTable(&TestModel{})

	err = db.AutoMigrate(&TestModel{})
	assert.NoError(t, err)

	return db
}

func TestGenerics(t *testing.T) {
	db := setupTestDB(t)

	// Test Create
	item := &TestModel{Name: "Alice", Age: 30}
	err := Create(db, item)
	assert.NoError(t, err)
	assert.NotZero(t, item.ID)

	// Test GetByID
	fetched, err := GetByID[TestModel](db, item.ID)
	assert.NoError(t, err)
	assert.Equal(t, item.Name, fetched.Name)

	// Test GetAll
	item2 := &TestModel{Name: "Bob", Age: 25}
	err = Create(db, item2)
	assert.NoError(t, err)

	all, err := GetAll[TestModel](db)
	assert.NoError(t, err)
	assert.Len(t, all, 2)

	// Test Update
	item.Age = 31
	err = Update(db, item)
	assert.NoError(t, err)

	fetchedUpdated, err := GetByID[TestModel](db, item.ID)
	assert.NoError(t, err)
	assert.Equal(t, 31, fetchedUpdated.Age)

	// Test Delete
	err = Delete(db, item)
	assert.NoError(t, err)

	allAfterDelete, err := GetAll[TestModel](db)
	assert.NoError(t, err)
	assert.Len(t, allAfterDelete, 1)
	assert.Equal(t, "Bob", allAfterDelete[0].Name)
}

func TestGetPage(t *testing.T) {
	db := setupTestDB(t)

	// Seed data
	for i := 0; i < 10; i++ {
		Create(db, &TestModel{Name: "User", Age: i})
	}

	// Test Pagination
	// Page 1: 0-4
	items, total, err := GetPage[TestModel](db, 0, 5, "age asc", nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
	assert.Len(t, items, 5)
	assert.Equal(t, 0, items[0].Age)
	assert.Equal(t, 4, items[4].Age)

	// Page 2: 5-9
	items, total, err = GetPage[TestModel](db, 5, 5, "age asc", nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
	assert.Len(t, items, 5)
	assert.Equal(t, 5, items[0].Age)
	assert.Equal(t, 9, items[4].Age)
}

func TestGetPage_filtered(t *testing.T) {
	db := setupTestDB(t)

	// Seed data
	Create(db, &TestModel{Name: "Alice", Age: 30})
	Create(db, &TestModel{Name: "Bob", Age: 30})
	Create(db, &TestModel{Name: "Charlie", Age: 25})
	Create(db, &TestModel{Name: "Dave", Age: 40})

	// Test Filtered Pagination
	// Filter Age = 30
	items, total, err := GetPage[TestModel](db, 0, 10, "name asc", "age = ?", 30)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, items, 2)
	assert.Equal(t, "Alice", items[0].Name)
	assert.Equal(t, "Bob", items[1].Name)

	// Filter Name LIKE 'D%'
	items, total, err = GetPage[TestModel](db, 0, 10, "name asc", "name LIKE ?", "D%")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, items, 1)
	assert.Equal(t, "Dave", items[0].Name)
}
