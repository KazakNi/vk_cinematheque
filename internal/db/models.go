package db

type User struct {
	Id       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Is_admin bool   `db:"is_admin"`
}
