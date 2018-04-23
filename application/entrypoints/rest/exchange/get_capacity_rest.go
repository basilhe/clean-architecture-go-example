package exchange

import (
	usecase "github.com/basilhe/tdd/application/core/usecase/exchange"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
)

type CapacityDto interface {
	GetHasAdslCapacity() bool
	GetHasFibreCapacity() bool
}

type GetCapacityForExchangeEndpoint interface {
	GetCapacity(exchange string) (CapacityDto, error)
	GetCapacityHandler(w http.ResponseWriter, r *http.Request)
}

type getCapacityForExchangeEndpoint struct {
	usecase.IGetCapacityForExchangeUseCase
}

func NewGetCapacityForExchangeEndpoint(useCase usecase.IGetCapacityForExchangeUseCase) GetCapacityForExchangeEndpoint {
	return &getCapacityForExchangeEndpoint{useCase}
}

func (e *getCapacityForExchangeEndpoint) GetCapacity(exchange string) (CapacityDto, error) {
	capacity, err := e.IGetCapacityForExchangeUseCase.GetCapacity(exchange)

	if err != nil {
		return nil, err
	}

	return &capacityDto{capacity.HasAdslCapacity(), capacity.HasFibreCapacity()}, nil
}

func (e *getCapacityForExchangeEndpoint) GetCapacityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	capacity, err := e.GetCapacity(vars["exchange"])
	if err != nil && err == usecase.ErrorExchangeNotFound {
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	jsonCapacity, _ := json.Marshal(capacity)
	w.Write(jsonCapacity)
}

type capacityDto struct {
	HasAdslCapacity bool `json:"hasAdslCapacity"`
	HasFibreCapacity bool `json:"hasFibreCapacity"`
}

func (d *capacityDto) GetHasAdslCapacity() bool {
	return d.HasAdslCapacity
}

func (d *capacityDto) GetHasFibreCapacity() bool {
	return d.HasFibreCapacity
}