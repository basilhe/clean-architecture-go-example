package exchange

import (
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/basilhe/tdd/application/core/usecase/mock"
	"github.com/basilhe/tdd/application/core/usecase/exchange"
	"github.com/basilhe/tdd/application/core/entity"
)

const EXCHANGE_CODE = "exch1"
var mockIGetCapacityForExchangeUseCase *mock.MockIGetCapacityForExchangeUseCase

var endpointGetCapacityForExchange GetCapacityForExchangeEndpoint

func TestReturnsTheCapacityForAnExchange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenThereIsCapacityForAnExchange(ctrl)
	setupEndpoint()
	capacity, err := endpointGetCapacityForExchange.GetCapacity(EXCHANGE_CODE)

	assert.NoError(t, err)
	assert.True(t, capacity.GetHasAdslCapacity())
	assert.True(t, capacity.GetHasFibreCapacity())
}

func TestErrorWhenDeviceIsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	giveAnExchangeThatDoesNotExist(ctrl)
	setupEndpoint()

	_, err := endpointGetCapacityForExchange.GetCapacity(EXCHANGE_CODE)

	assert.Error(t, err)
}

func setupEndpoint() {
	endpointGetCapacityForExchange = NewGetCapacityForExchangeEndpoint(mockIGetCapacityForExchangeUseCase)
}

func givenThereIsCapacityForAnExchange(ctrl *gomock.Controller) {
	mockIGetCapacityForExchangeUseCase = mock.NewMockIGetCapacityForExchangeUseCase(ctrl)
	mockIGetCapacityForExchangeUseCase.EXPECT().GetCapacity(EXCHANGE_CODE).Return(&entity.Capacity{true, true}, nil)
}

func giveAnExchangeThatDoesNotExist(ctrl *gomock.Controller) {
	mockIGetCapacityForExchangeUseCase = mock.NewMockIGetCapacityForExchangeUseCase(ctrl)
	mockIGetCapacityForExchangeUseCase.EXPECT().GetCapacity(EXCHANGE_CODE).Return(nil, exchange.ErrorExchangeNotFound)
}