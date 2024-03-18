package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DBConnection *sqlx.DB

var GetListMovieStmt string = `Select f.name, f.description, f.release_date, f.rating, array_agg(a.name) as actors_list
FROM films f
JOIN castfilms cs
ON cs.film_id = f.id
JOIN actors a
ON a.id = cs.actor_id
GROUP BY f.name, f.description, f.release_date, f.rating
ORDER BY f.rating DESC`

var GetListMoviesSortByNameStmt string = `Select f.name, f.description, f.release_date, f.rating, array_agg(a.name) as actors_list
FROM films f
JOIN castfilms cs
ON cs.film_id = f.id
JOIN actors a
ON a.id = cs.actor_id
GROUP BY f.name, f.description, f.release_date, f.rating
ORDER BY f.name ASC`

var GetListMoviesSortByDateStmt string = `Select f.name, f.description, f.release_date, f.rating, array_agg(a.name) as actors_list
FROM films f
JOIN castfilms cs
ON cs.film_id = f.id
JOIN actors a
ON a.id = cs.actor_id
GROUP BY f.name, f.description, f.release_date, f.rating
ORDER BY f.release_date ASC`

var GetListMoviesSortByRatingStmt string = `Select f.name, f.description, f.release_date, f.rating, array_agg(a.name) as actors_list
FROM films f
JOIN castfilms cs
ON cs.film_id = f.id
JOIN actors a
ON a.id = cs.actor_id
GROUP BY f.name, f.description, f.release_date, f.rating
ORDER BY f.rating ASC`

var GetListMoviesSearchByActorStmt string = `Select f.name, f.description, f.release_date, f.rating, array_agg(a.name) as actors_list
FROM films f
JOIN castfilms cs
ON cs.film_id = f.id
JOIN actors a
ON a.id = cs.actor_id
WHERE a.name ILIKE '%' || $1 || '%'
GROUP BY f.name, f.description, f.release_date, f.rating
ORDER BY f.rating DESC
`

var GetListMoviesSearchByMovieStmt string = `Select f.name, f.description, f.release_date, f.rating, array_agg(a.name) as actors_list
FROM films f
JOIN castfilms cs
ON cs.film_id = f.id
JOIN actors a
ON a.id = cs.actor_id
WHERE f.name ILIKE '%' || $1 || '%'
GROUP BY f.name, f.description, f.release_date, f.rating
ORDER BY f.rating DESC`

var GetListActorsStmt string = `Select a.name, a.sex, a.birth_date, array_agg(f.name) as actor_films
FROM actors a
JOIN castfilms cs
ON cs.actor_id = a.id
JOIN films f
ON f.id = cs.film_id
GROUP BY a.name, a.sex, a.birth_date
`

func NewDBConnection() (*sqlx.DB, error) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parent := filepath.Dir(wd)

	err = godotenv.Load(parent + "./.env")
	if err != nil {
		panic("Error loading .env file")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DB_NAME")
	driver := os.Getenv("DRIVER")

	connUrl := fmt.Sprintf("host=%s port=%v user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Connect(driver, connUrl)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to DB")
	return db, nil
}

func ExecMigration(db *sqlx.DB, path string) error {

	query, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
	return nil
}
