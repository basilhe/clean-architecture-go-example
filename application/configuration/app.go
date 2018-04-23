package main

import (
	"github.com/gorilla/mux"
	rest "github.com/basilhe/tdd/application/entrypoints/rest/exchange"
	usecase "github.com/basilhe/tdd/application/core/usecase/exchange"
	provider "github.com/basilhe/tdd/application/dataproviders/exchange"
	"github.com/basilhe/tdd/application/core/entity"
	"net/http"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

type dataprovider struct {

}

func (d *dataprovider) GetAvailablePortsOfAllDevicesInExchange(exchange string) ([]*entity.BroadbandAccessDevice, error) {
	return []*entity.BroadbandAccessDevice{
		{Hostname: "google.com", SerialNumber: "google", AvailablePorts: 2, DeviceType: entity.ADSL},
		{Hostname: "baidu.com", SerialNumber: "baidu", AvailablePorts: 2, DeviceType: entity.FIBRE},
		{Hostname: "163.com", SerialNumber: "163", AvailablePorts: 5, DeviceType: entity.ADSL},
	}, nil
}

func main() {
	r := mux.NewRouter()

	db := prepareDB()

	dpExchangeExist := provider.NewExchangeDBDataProvider(db)
	dp := &dataprovider{}
	usecaseGetCapacityForExchange := usecase.NewGetCapacityForExchangeUseCase(dp, dpExchangeExist)
	endpoint := rest.NewGetCapacityForExchangeEndpoint(usecaseGetCapacityForExchange)
	r.HandleFunc("/exchange/{exchange}/capacity", endpoint.GetCapacityHandler)

	srv := &http.Server{
		Addr:    "127.0.0.1:8000",
		Handler: r,
	}

	log.Fatal(srv.ListenAndServe())
}

func prepareDB() *sql.DB {
	db, e := sql.Open("sqlite3", ":memory:")
	if e != nil {
		log.Fatal(e)
	}
	tx, err := db.Begin()
	_, err = db.Exec(`CREATE TABLE exchange (code VARCHAR(50) PRIMARY KEY, count INT);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO exchange VALUES("exch1", 1);`)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	return db
}
