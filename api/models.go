package api

import (
	"cinematheque/internal/db"
	"database/sql"
	"log"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) IsUserExists(db *sqlx.DB) (exists bool, err error) {
	var id int
	err = db.Get(&id, "SELECT id FROM users WHERE email = $1", user.Email)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		log.Printf("error while querying: %s", err)
		return false, err
	} else {
		return true, nil
	}
}

func (user *User) CreateUser(db *sqlx.DB) (user_id int, err error) {
	var id int
	row := db.QueryRow("INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id", user.Email, user.Password)
	err = row.Scan(&id)
	if err != nil {
		log.Printf("error while insert values: %s", err)
		return id, err
	}
	return id, nil

}

func (user *User) GetUserByEmail(email string, database *sqlx.DB) (dbuser db.User, err error) {
	var db_user db.User
	err = database.Get(&db_user, "SELECT id, email, password FROM users WHERE email = $1", email)
	if err == sql.ErrNoRows {
		log.Println("No email in db!")
		return db.User{}, nil
	} else if err != nil {
		log.Printf("error while querying user: %s", err)
		return db_user, err
	} else {
		return db_user, nil
	}
}

func (user *User) CheckPassword(providedPassword string, db_password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(db_password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

type UserId struct {
	Id int `json:"id" db:"id"`
}

type CustomClaims struct {
	Is_admin bool `json:"is_admin"`
	jwt.RegisteredClaims
}
