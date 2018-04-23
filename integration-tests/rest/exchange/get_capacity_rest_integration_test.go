package exchange

import (
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/basilhe/tdd/application/core/usecase/mock"
	rest "github.com/basilhe/tdd/application/entrypoints/rest/exchange"
	usecase "github.com/basilhe/tdd/application/core/usecase/exchange"

	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"github.com/basilhe/tdd/application/core/entity"
)

const EXCHANGE_CODE = "exch1"
var mockIGetCapacityForExchangeUseCase *mock.MockIGetCapacityForExchangeUseCase

var endpointGetCapacityForExchange rest.GetCapacityForExchangeEndpoint

var responseStatusCode int
var responseContent string

func TestReturnsCapacityForAnExchange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenThereIsCapacityForTheExchange(ctrl)
	setupEndpoint(ctrl)

	whenTheCapacityIsRetrieved(t)

	thenTheCapacityIsRetruned(t)
}

func TestReturns404WhenTheExchangeIsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenTheExchangeDoesNotExist(ctrl)
	setupEndpoint(ctrl)

	whenTheCapacityIsRetrieved(t)

	thenItReturns404(t)
}

func whenTheCapacityIsRetrieved(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/exchange/{exchange}/capacity", endpointGetCapacityForExchange.GetCapacityHandler)

	req, err := http.NewRequest("GET", "/exchange/exch1/capacity", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	responseStatusCode = rr.Code

	responseContent = rr.Body.String()
}

func setupEndpoint(ctrl *gomock.Controller) {
	endpointGetCapacityForExchange = rest.NewGetCapacityForExchangeEndpoint(mockIGetCapacityForExchangeUseCase)
}

func thenItReturns404(t *testing.T) {
	assert.Equal(t, http.StatusNotFound, responseStatusCode)
}

func givenTheExchangeDoesNotExist(ctrl *gomock.Controller) {
	mockIGetCapacityForExchangeUseCase = mock.NewMockIGetCapacityForExchangeUseCase(ctrl)
	mockIGetCapacityForExchangeUseCase.EXPECT().GetCapacity(EXCHANGE_CODE).Return(nil, usecase.ErrorExchangeNotFound)
}

func thenTheCapacityIsRetruned(t *testing.T) {
	assert.Equal(t, http.StatusOK, responseStatusCode)
	assert.Equal(t, responseContent,"{\"hasAdslCapacity\":true,\"hasFibreCapacity\":false}")
}

func givenThereIsCapacityForTheExchange(ctrl *gomock.Controller) {
	mockIGetCapacityForExchangeUseCase = mock.NewMockIGetCapacityForExchangeUseCase(ctrl)
	mockIGetCapacityForExchangeUseCase.EXPECT().GetCapacity(EXCHANGE_CODE).Return(&entity.Capacity{true, false}, nil)
}