package testsetup

import (
	"context"
	"database/sql"
	"testing"

	"github.com/volatiletech/null"

	"github.com/labstack/echo/v4"

	"github.com/nmrshll/kata-payments-api/api"
	models "github.com/nmrshll/kata-payments-api/generated-models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"

	// load drivers
	_ "github.com/lib/pq"
	_ "github.com/volatiletech/sqlboiler/queries/qm"
)

var bg = context.Background()

type TestSetup struct {
	FnCleanup  func()
	DB         *sql.DB
	EchoServer *echo.Echo
	Data       *TestData
}

func New(t *testing.T) (*TestSetup, *TestData) {
	db, err := sql.Open("postgres", "postgres://dbuser:dbpass@db:5432/dbname?sslmode=disable")
	assert.NoError(t, err)

	// If you don't want to pass in db to all generated methods
	// you can use boil.SetDB to set it globally, and then use
	// the G variant methods like so (--add-global-variants to enable)
	boil.SetDB(db)

	e := api.NewServer()

	fnCleanup := func() {
		_, err := queries.Raw(`TRUNCATE payments, currencies, parties`).Exec(boil.GetDB())
		assert.NoError(t, err)
	}
	testData := insertTestData(t)

	return &TestSetup{fnCleanup, db, e, testData}, testData
}

type TestData struct {
	Currency1 *models.Currency
	Party1    *models.Party
	Party2    *models.Party
}

func insertTestData(t *testing.T) *TestData {
	currency1 := &models.Currency{
		Name: null.StringFrom("Euro"), Symbol: null.StringFrom("EUR"),
	}
	err := currency1.InsertG(bg, boil.Infer())
	assert.NoError(t, err)

	party1 := &models.Party{
		AccountName:   null.StringFrom("party1"),
		AccountNumber: null.StringFrom("df654gg7"),
	}
	err = party1.InsertG(bg, boil.Infer())
	assert.NoError(t, err)

	party2 := &models.Party{
		AccountName:   null.StringFrom("party2"),
		AccountNumber: null.StringFrom("b7v85g6"),
	}
	err = party2.InsertG(bg, boil.Infer())
	assert.NoError(t, err)

	return &TestData{
		Currency1: currency1,
		Party1:    party1,
		Party2:    party2,
	}
}
