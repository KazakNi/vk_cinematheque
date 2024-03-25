package utils

import (
	"os"

	"github.com/joho/godotenv"
)

var dev = false

func GetTestMigrPathDown() string {
	if dev {
		return "./internal/migrations/test_migration_down.sql"
	} else {
		return "../internal/migrations/test_migration_down.sql"
	}
}

func GetTestMigrPathUp() string {
	if dev {
		return "./internal/migrations/test_migration.sql"
	} else {
		return "../internal/migrations/test_migration.sql"
	}
}

func GetMigrPath() string {
	if dev {
		return "./internal/migrations/migrations.sql"
	} else {
		return "../internal/migrations/migrations.sql"
	}
}

func GetStaticPath() string {
	if dev {
		return "./api/static/redoc.html"
	} else {
		return "../api/static/redoc.html"
	}
}

func GetStaticRoot() string {

	if dev {
		return "./api/static"
	} else {
		return "../api/static"
	}
}

func GetEnv() (hostDB, portDB, userDB, passwordDB, dbnameDB, driverDB string) {

	var host, port, user, password, dbname, driver string
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../example.env")
		if err != nil {
			panic("error while loading .env file")
		}

		host = os.Getenv("HOST_LOCAL")
	} else {
		host = os.Getenv("HOST")
	}

	port = os.Getenv("PORT")
	user = os.Getenv("USER")
	password = os.Getenv("PASSWORD")
	dbname = os.Getenv("DB_NAME")
	driver = os.Getenv("DRIVER")

	return host, port, user, password, dbname, driver
}
