package exchange

import (
	"database/sql"
	usecase "github.com/basilhe/tdd/application/core/usecase/exchange"
	"fmt"
)

type exchangeDataProvider struct {
	db *sql.DB
}

func (provider *exchangeDataProvider) DoesExchangeExist(exchange string) bool {
	var count int

	row := provider.db.QueryRow(`SELECT count(1) FROM exchange WHERE code = ?`, exchange)

	err := row.Scan(&count)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return false;
	}
	return count >= 1
}

func NewExchangeDBDataProvider(db *sql.DB) usecase.IDoesExchangeExist {
	return &exchangeDataProvider{db}
}