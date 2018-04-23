package exchange

import (
	"testing"
	"github.com/stretchr/testify/assert"
	usecase "github.com/basilhe/tdd/application/core/usecase/exchange"
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const EXCHANGE_CODE = "exch1"
var exchangeDBDataProvider usecase.IDoesExchangeExist
var db *sql.DB
var mock sqlmock.Sqlmock

func TestFalseWhenExchangeDoesExists(t *testing.T) {
	db = givenExchangeExist(t, false)
	exchangeDBDataProvider = NewExchangeDBDataProvider(db)

	doesExchangeExist := exchangeDBDataProvider.DoesExchangeExist(EXCHANGE_CODE)

	assert.False(t, doesExchangeExist)
}

func TestTrueWhenExchangeDoesExist(t *testing.T) {

	db = givenExchangeExist(t, true)
	exchangeDBDataProvider = NewExchangeDBDataProvider(db)

	doesExchangeExist := exchangeDBDataProvider.DoesExchangeExist(EXCHANGE_CODE)

	assert.True(t, doesExchangeExist)
}

func givenExchangeExist(t *testing.T, exist bool) *sql.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	var count int
	if exist {
		count = 1
	}
	mock.ExpectQuery("SELECT count\\(1\\) FROM exchange WHERE code = \\?").WithArgs("exch1").WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(count))
	return db
}
