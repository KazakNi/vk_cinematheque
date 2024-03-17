package api

import (
	"cinematheque/internal/db"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

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
	err = database.Get(&db_user, "SELECT id, email, password, is_admin FROM users WHERE email = $1", email)
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

type CreatedId struct {
	Id int `json:"id" db:"id"`
}

type CustomClaims struct {
	Is_admin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

type PostFilm struct {
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	Release_date time.Time `json:"release_date" db:"release_date"`
	Rating       int       `json:"rating" db:"rating"`
	Actors_list  []int     `json:"actors_list"` // actors ids for MtM relation support
}

type Film struct {
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	Release_date time.Time `json:"release_date" db:"release_date"`
	Rating       int       `json:"rating" db:"rating"`
}

func (film *Film) Create(db *sqlx.DB) (film_id int, err error) {
	var id int
	row := db.QueryRow("INSERT INTO films (name, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING id", film.Name, film.Description, film.Release_date, film.Rating)
	err = row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (film *Film) Delete(id int, db *sqlx.DB) error {
	var ID int
	err := db.QueryRow("DELETE FROM films WHERE id = $1 RETURNING id", id).Scan(&ID)
	if err == sql.ErrNoRows {
		return err
	}
	return nil
}

func (film *Film) Update(id int, actors_ids []int, db *sqlx.DB) error {
	var ID int
	err := db.QueryRow("UPDATE films SET name = COALESCE($1, name), description = COALESCE($2, description), release_date = COALESCE($3, release_date)  rating = COALESCE($4, rating) WHERE id = $5 RETURNING id", film.Name, film.Description, film.Release_date, film.Rating, id).Scan(&ID)
	if err == sql.ErrNoRows {
		return err
	}

	db.MustExec("DELETE FROM castfilms WHERE film_id = $1", id)

	err = film.InsertCast(id, actors_ids, db)
	if err != nil {
		return err
	}
	return nil
}

func (film *Film) GetFilmById(id int, db *sqlx.DB) (db_film Film, err error) {
	var f Film
	err = db.Get(&f, "SELECT name, description, release_date, rating FROM films WHERE id = $1", id)
	if err == sql.ErrNoRows {
		log.Println("No film in db!")
		return Film{}, err
	} else if err != nil {
		log.Printf("error while querying film: %s", err)
		return f, err
	} else {
		return f, nil
	}
}

func (film *Film) InsertCast(film_id int, actors_ids []int, db *sqlx.DB) (err error) {

	castMaps := []map[string]interface{}{}
	for _, id := range actors_ids {
		var m = make(map[string]interface{})
		m["actor_id"] = id
		m["film_id"] = film_id

		castMaps = append(castMaps, m)
	}

	_, err = db.NamedExec(`INSERT INTO castfilms (actor_id, film_id)
        VALUES (:actor_id, :film_id)`, castMaps)
	if err != nil {
		return err
	}
	return nil
}

type Actor struct {
	Name       string    `json:"name" db:"name"`
	Sex        string    `json:"sex" db:"sex"`
	Birth_date BirthDate `json:"birth_date" db:"birth_date"`
}

func (actor *Actor) Create(db *sqlx.DB) (actor_id int, err error) {
	var id int
	row := db.QueryRow("INSERT INTO actors (name, sex, birth_date) VALUES ($1, $2, $3) RETURNING id", actor.Name, actor.Sex, actor.Birth_date.Time)
	err = row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (actor *Actor) Delete(id int, db *sqlx.DB) error {
	var ID int
	err := db.QueryRow("DELETE FROM actors WHERE id = $1 RETURNING id", id).Scan(&ID)
	if err == sql.ErrNoRows {
		return err
	}
	return nil
}

func (actor *Actor) Update(id int, db *sqlx.DB) error {
	var ID int
	err := db.QueryRow("UPDATE actors SET name = COALESCE($1, name), sex = COALESCE($2, sex), birth_date = COALESCE($3, birth_date)  WHERE id = $4 RETURNING id", actor.Name, actor.Sex, actor.Birth_date.Time, id).Scan(&ID)
	if err == sql.ErrNoRows {
		return err
	}
	return nil
}

func (actor *Actor) GetActorById(id int, db *sqlx.DB) (db_actor Actor, err error) {
	var a Actor
	err = db.Get(&a, "SELECT name, sex, birth_date FROM actors WHERE id = $1", id)
	if err == sql.ErrNoRows {
		log.Println("No actor in db!")
		return Actor{}, err
	} else if err != nil {
		log.Printf("error while querying user: %s", err)
		return a, err
	} else {
		return a, nil
	}
}

type BirthDate struct {
	time.Time
}

func (date *BirthDate) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, _ := time.Parse("2006-01-02", s)
	date.Time = t
	return nil
}

func (date *BirthDate) Scan(src interface{}) error {
	var value string
	switch src.(type) {
	case time.Time:
		value = fmt.Sprintf("%s", src)
	default:
		return errors.New("invalid type for BirthDate")
	}
	t, _ := time.Parse("2006-01-02", value)
	*date = BirthDate{t}
	return nil
}

func (date *BirthDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(date.Time.Format("2006-01-02"))
}
