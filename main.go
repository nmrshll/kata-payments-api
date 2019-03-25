package main

import (
	"database/sql"
	"fmt"
	"strconv"

	// Import this so we don't have to use qm.Limit etc.
	_ "github.com/volatiletech/sqlboiler/queries/qm"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/nmrshll/kata-payments-api/api"
)

const PORT = 8080

func main() {
	db, err := sql.Open("postgres", "postgres://dbuser:dbpass@db:5432/dbname?sslmode=disable")
	if err != nil {
		panic(err)
	}

	// If you don't want to pass in db to all generated methods
	// you can use boil.SetDB to set it globally, and then use
	// the G variant methods like so (--add-global-variants to enable)
	boil.SetDB(db)

	e := api.NewServer()

	// Start server
	fmt.Printf("Listening on port %s...\n", strconv.Itoa(PORT))
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))
}
