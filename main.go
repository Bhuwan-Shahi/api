package main

// Model for Journal - file
type Journal struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	User  *User  `json:"user"`
}

type User struct {
	Id   int64  `json:"userId"`
	Name string `json:"name"`
}

var journal []Journal

// Middleware
func isEmpty(J *Journal) bool {
	return J.Title == ""
}

func isempty(U *User) bool {
	return U.Name == ""
}
func main() {

}
