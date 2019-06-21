package models

// User is login information
type User struct {
	Username string `db:"username"`
	Secret   []byte `db:"secret"`
}
