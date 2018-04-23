package exchange

import (
	"github.com/stretchr/testify/assert"
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/basilhe/tdd/application/core/usecase/mock"
	"github.com/basilhe/tdd/application/core/entity"
)

const NO_PORTS = 0
const EXCHANGE_CODE = "exchange_code"

var mockIGetAvailablePortsOfAllDevicesInExchange *mock.MockIGetAvailablePortsOfAllDevicesInExchange
var mockIDoesExchangeExist *mock.MockIDoesExchangeExist
var useCaseGetCapacityForExchange IGetCapacityForExchangeUseCase

func TestNoCapacityIfDevicesHaveNoPorts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenDevices(ctrl, device(entity.FIBRE, NO_PORTS), device(entity.ADSL, NO_PORTS))
	givenExchange(ctrl, true)
	setupUseCase()

	capacity, err := useCaseGetCapacityForExchange.GetCapacity(EXCHANGE_CODE)

	assert.Nil(t, err, "Error should not occur")
	assert.False(t, capacity.HasAdslCapacity(), "Adsl has no capacity")
	assert.False(t, capacity.HasFibreCapacity(), "Fibre has no capacity")
}

func TestAdslHasNoCapacityIfAvailablePortsIsLessThan5(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenDevices(ctrl, device(entity.ADSL, 1), device(entity.ADSL, 3))
	givenExchange(ctrl, true)
	setupUseCase()

	capacity, err := useCaseGetCapacityForExchange.GetCapacity(EXCHANGE_CODE)

	assert.Nil(t, err, "Error should not occur")
	assert.False(t, capacity.HasAdslCapacity())
}

func TestFibreHasNoCapacityIfAvailablePortsIsLessThan5(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenDevices(ctrl, device(entity.FIBRE, 1), device(entity.FIBRE, 3))
	givenExchange(ctrl, true)
	setupUseCase()

	capacity, err := useCaseGetCapacityForExchange.GetCapacity(EXCHANGE_CODE)

	assert.Nil(t, err, "Error should not occur")
	assert.False(t, capacity.HasFibreCapacity())
}

func TestAdslHasCapacityIfAvailablePortsIsEqualTo5(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenDevices(ctrl, device(entity.ADSL, 2), device(entity.ADSL, 3))
	givenExchange(ctrl, true)
	setupUseCase()

	capacity, err := useCaseGetCapacityForExchange.GetCapacity(EXCHANGE_CODE)

	assert.Nil(t, err, "Error should not occur")
	assert.True(t, capacity.HasAdslCapacity())
}

func TestFibreHasCapacityIfAvailablePortsIsEqualTo5(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenDevices(ctrl, device(entity.FIBRE, 2), device(entity.FIBRE, 3))
	givenExchange(ctrl, true)
	setupUseCase()

	capacity, err := useCaseGetCapacityForExchange.GetCapacity(EXCHANGE_CODE)

	assert.Nil(t, err, "Error should not occur")
	assert.True(t, capacity.HasFibreCapacity())
}

func TestErrorWhenExchangeDoesNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenExchange(ctrl, false)
	setupUseCase()
	_, err := useCaseGetCapacityForExchange.GetCapacity(EXCHANGE_CODE)

	assert.NotNil(t, err, "Error should not occur")
}

func setupUseCase() {
	useCaseGetCapacityForExchange = NewGetCapacityForExchangeUseCase(mockIGetAvailablePortsOfAllDevicesInExchange, mockIDoesExchangeExist)
}

func givenDevices(ctrl *gomock.Controller, devices ...*entity.BroadbandAccessDevice) {
	mockIGetAvailablePortsOfAllDevicesInExchange = mock.NewMockIGetAvailablePortsOfAllDevicesInExchange(ctrl)
	mockIGetAvailablePortsOfAllDevicesInExchange.EXPECT().GetAvailablePortsOfAllDevicesInExchange(EXCHANGE_CODE).Return(devices, nil)

}

func givenExchange(ctrl *gomock.Controller, exist bool) {
	mockIDoesExchangeExist = mock.NewMockIDoesExchangeExist(ctrl)
	mockIDoesExchangeExist.EXPECT().DoesExchangeExist(EXCHANGE_CODE).Return(exist)
}

func device(deviceType entity.DeviceType, port int) *entity.BroadbandAccessDevice {
	return &entity.BroadbandAccessDevice{
		"hostname", "serial", deviceType, port,
	}
}
