package usecases

import (
	"rest-api/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	db := NewDatabase()
	id := db.Add(domain.Book{
		Author: "Marin Preda",
		Title:  "Cel mai iubit dintre pamanteni",
		Genre:  "Novel",
		Year:   1980,
	})

	book, err := db.GetByID(id)
	assert.NoError(t, err, "Should be nil")
	assert.Equal(t, book.Author, "Marin Preda")
	assert.Equal(t, book.Title, "Cel mai iubit dintre pamanteni")
	assert.Equal(t, book.Genre, "Novel")
	assert.Equal(t, book.Year, 1980)

	books := db.GetAll(make(map[string]interface{}))
	assert.Equal(t, len(books), 1)
	assert.Equal(t, books[0], book)

	db.Update(id, domain.Book{
		Author: "Marin Preda",
		Title:  "Morometii",
		Genre:  "Novel",
		Year:   1955,
	})

	book, err = db.GetByID(id)
	assert.NoError(t, err, "Should be nil")
	assert.Equal(t, book.Author, "Marin Preda")
	assert.Equal(t, book.Title, "Morometii")
	assert.Equal(t, book.Genre, "Novel")
	assert.Equal(t, book.Year, 1955)

	err = db.RemoveByID(id)
	assert.NoError(t, err, "Should be nil")

	_, err = db.GetByID(id)
	assert.Error(t, err, "Should be not found")
}

func TestFilters(t *testing.T) {
	db := NewDatabase()
	db.Add(domain.Book{
		Author: "Marin Preda",
		Title:  "Cel mai iubit dintre pamanteni",
		Genre:  "Novel",
		Year:   1980,
	})

	db.Add(domain.Book{
		Author: "Marin Preda",
		Title:  "Morometii",
		Genre:  "Novel",
		Year:   1955,
	})

	db.Add(domain.Book{
		Author: "Marin Preda",
		Title:  "Delirul",
		Genre:  "Novel",
		Year:   1975,
	})

	db.Add(domain.Book{
		Author: "Mircea Cartarescu",
		Title:  "Theodoros",
		Genre:  "Novel",
		Year:   2022,
	})

	filtersAuthor := map[string]interface{}{
		"author": "Marin Preda",
	}

	books := db.GetAll(filtersAuthor)
	assert.Equal(t, len(books), 3)

	filtersYear := map[string]interface{}{
		"year": 1980,
	}

	books = db.GetAll(filtersYear)
	assert.Equal(t, len(books), 1)

	filtersAuthorAndYear := map[string]interface{}{
		"author": "Marin Preda",
		"year":   1980,
	}

	books = db.GetAll(filtersAuthorAndYear)
	assert.Equal(t, len(books), 1)
}
