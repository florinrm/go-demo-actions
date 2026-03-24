package usecases

import "rest-api/domain"

//go:generate mockery --name Database --output ../mocks --filename=mock_database.go

// Database has a set of operations for adding in a key-value database
type Database interface {
	Add(book domain.Book) string
	Update(hash string, book domain.Book)
	RemoveByID(id string) error
	GetByID(id string) (domain.Book, error)
	GetAll(filters map[string]interface{}) []domain.Book
}
