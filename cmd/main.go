package main

import (
	"cinematheque/api"
	"cinematheque/internal/db"
	"cinematheque/internal/utils"
)

func main() {
	migration_path := utils.GetMigrPath()
	db.DBConnection, _ = db.NewDBConnection()
	db.ExecMigration(db.DBConnection, migration_path)
	api.SetRoutes()
}
