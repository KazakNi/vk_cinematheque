package db

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	Id       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Is_admin bool   `db:"is_admin"`
}

type Film struct {
	Name         string         `json:"name" db:"name"`
	Description  string         `json:"description" db:"description"`
	Release_date time.Time      `json:"release_date" db:"release_date"`
	Rating       int            `json:"rating" db:"rating"`
	Actors_list  pq.StringArray `json:"actors_list" db:"actors_list"`
}
