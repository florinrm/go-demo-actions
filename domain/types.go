package domain

// Book encapsulated data about a published book
type Book struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Genre  string `json:"genre"`
	Year   int    `json:"year"`
}

// Response is the API response for GET, POST, PUT, PATCH
type Response struct {
	ID string `json:"id"`
	Book
}
