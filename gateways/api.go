package gateways

import (
	"encoding/json"
	"net/http"
	"rest-api/domain"
	"rest-api/usecases"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

// BooksAPI represent the API used for managing books
type BooksAPI struct {
	database usecases.Database
	logger   *zerolog.Logger
	router   *mux.Router
}

// NewBooksAPI creates new instance of BooksAPI
func NewBooksAPI(database usecases.Database, logger *zerolog.Logger) *BooksAPI {
	router := mux.NewRouter()
	return &BooksAPI{
		database: database,
		logger:   logger,
		router:   router,
	}
}

func (api *BooksAPI) GetRouter() *mux.Router {
	return api.router
}

func (api *BooksAPI) Handle() {
	api.router.HandleFunc("/entities/books", api.AddBook).Methods(http.MethodPost)
	api.router.HandleFunc("/queries/books", api.GetBooks).Methods(http.MethodGet)
	api.router.HandleFunc("/entities/books/{id}", api.GetBook).Methods(http.MethodGet)
	api.router.HandleFunc("/entities/books/{id}", api.DeleteBook).Methods(http.MethodDelete)
	api.router.HandleFunc("/entities/books/{id}", api.PutBook).Methods(http.MethodPut)
	api.router.HandleFunc("/entities/books/{id}", api.PatchBook).Methods(http.MethodPatch)
}

func (api *BooksAPI) Serve() error {
	err := http.ListenAndServe(":8080", api.router)
	if err != nil {
		return err
	}

	return nil
}

// GetBooks used for handling GET endpoint for getting all the books
func (api *BooksAPI) GetBooks(w http.ResponseWriter, r *http.Request) {
	filters, err := api.getQueryParametersAllBooks(r)
	if err != nil {
		api.logger.Error().Err(err).Msg("invalid query parameters")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	books := api.database.GetAll(filters)
	err = encoder.Encode(books)
	if err != nil {
		api.logger.Error().Err(err).Msg("error encoding books")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (api *BooksAPI) getQueryParametersAllBooks(r *http.Request) (map[string]interface{}, error) {
	filters := make(map[string]interface{})

	title := r.URL.Query().Get("title")
	if len(title) != 0 {
		filters["title"] = title
	}

	author := r.URL.Query().Get("author")
	if len(title) != 0 {
		filters["author"] = author
	}

	genre := r.URL.Query().Get("genre")
	if len(title) != 0 {
		filters["genre"] = genre
	}

	year := r.URL.Query().Get("year")
	if len(title) != 0 {
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			return nil, err
		}
		filters["year"] = yearInt
	}

	return filters, nil
}

// GetBook used for handling GET endpoint for getting a book by ID
func (api *BooksAPI) GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, ok := params["id"]
	if !ok {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	book, err := api.database.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = encoder.Encode(domain.Response{
		ID:   id,
		Book: book,
	})
	if err != nil {
		api.logger.Error().Err(err).Msg("error encoding books")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AddBook used for handling POST endpoint for adding books
func (api *BooksAPI) AddBook(w http.ResponseWriter, r *http.Request) {
	var book domain.Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if err != nil {
		api.logger.Error().Err(err).Msg("error decoding book")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := api.database.Add(book)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(domain.Response{
		ID:   id,
		Book: book,
	})
	if err != nil {
		api.logger.Error().Err(err).Msg("error encoding books")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteBook used for handling DELETE endpoint for deleting a book by ID
func (api *BooksAPI) DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, ok := params["id"]
	if !ok {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	err := api.database.RemoveByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(id)
	if err != nil {
		api.logger.Error().Err(err).Msg("error encoding books")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// PutBook used for handling PUT endpoint for updating a book by ID
func (api *BooksAPI) PutBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, ok := params["id"]
	if !ok {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	var book domain.Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if err != nil {
		api.logger.Error().Err(err).Msg("error decoding book")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.database.Update(id, book)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(domain.Response{
		ID:   id,
		Book: book,
	})
	if err != nil {
		api.logger.Error().Err(err).Msg("error encoding books")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// PatchBook used for handling PATCH endpoint for updating a book by ID
func (api *BooksAPI) PatchBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, ok := params["id"]
	if !ok {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	var book domain.Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if err != nil {
		api.logger.Error().Err(err).Msg("error decoding book")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	originalBook, err := api.database.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if len(book.Title) != 0 {
		originalBook.Title = book.Title
	}

	if len(book.Author) != 0 {
		originalBook.Author = book.Author
	}

	if len(book.Genre) != 0 {
		originalBook.Genre = book.Genre
	}

	if book.Year != 0 {
		originalBook.Year = book.Year
	}

	api.database.Update(id, originalBook)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(domain.Response{
		ID:   id,
		Book: originalBook,
	})
	if err != nil {
		api.logger.Error().Err(err).Msg("error encoding books")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
