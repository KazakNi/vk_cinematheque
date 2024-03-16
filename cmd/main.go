package main

import (
	"cinematheque/api"
	"cinematheque/internal/db"
)

func main() {
	db.DBConnection, _ = db.NewDBConnection()
	//db.ExecMigration(db.DBConnection, "../internal/migrations/migrations.sql")
	api.SetRoutes()

}
