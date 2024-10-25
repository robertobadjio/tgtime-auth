package model

// User ???
type User struct {
	ID       string `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
	Role     string `db:"role"`
}
