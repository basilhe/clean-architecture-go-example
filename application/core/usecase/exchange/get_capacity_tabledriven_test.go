package exchange

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/basilhe/tdd/application/core/entity"
)

var getCapacityTests = []struct {
	testName                 string
	exchangeExist            bool
	deviceTypeOfFirstDevice  entity.DeviceType
	portsOfFirstDevice       int
	deviceTypeOfSecondDevice entity.DeviceType
	portsOfSecondDevice      int
	expectError              bool
	expectHasAdslCapacity    bool
	expectHasFibreCapacity   bool
}{
	{"noCapacityIfDevicesHasNoPorts", true, entity.FIBRE, NO_PORTS, entity.ADSL, NO_PORTS, false, false, false},
	{"adslHasNoCapacityIfAvailablePortsIsLessThan5", true, entity.ADSL, 1, entity.ADSL, 3, false, false, false},
	{"fibrehasNoCapacityIfAvailablePortsIsLessThan5", true, entity.FIBRE, 1, entity.FIBRE, 3, false, false, false},
	{"adslHasCapacityIfAvailablePortsIsEqualTo5", true, entity.ADSL, 2, entity.ADSL, 3, false, true, false},
	{"fibreHasCapacityIfAvailablePortsIsEqualTo5", true, entity.FIBRE, 2, entity.FIBRE, 3, false, false, true},
	{"errorWhenExchangeDoesNotExist", false, entity.ADSL, NO_PORTS, entity.FIBRE, NO_PORTS, true, false, false},
}

func TestGetCapacityAllInOne(t *testing.T) {
	for _, test := range getCapacityTests {
		ctrl := gomock.NewController(t)

		givenExchange(ctrl, test.exchangeExist)
		if test.exchangeExist {
			givenDevices(ctrl, device(test.deviceTypeOfFirstDevice, test.portsOfFirstDevice), device(test.deviceTypeOfSecondDevice, test.portsOfSecondDevice))
		}
		setupUseCase()

		capacity, err := useCaseGetCapacityForExchange.GetCapacity(EXCHANGE_CODE)

		if test.expectError {
			assert.NotNil(t, err, "Error should occur")
		} else {
			assert.Nil(t, err, "Error should not occur")
			assert.Equal(t, test.expectHasAdslCapacity, capacity.HasAdslCapacity(), "Test %s expects Adsl capacity equal to %v", test.testName, test.expectHasAdslCapacity)
			assert.Equal(t, test.expectHasFibreCapacity, capacity.HasFibreCapacity(), "Test %s expects Fibre capacity equal to %v", test.testName, test.expectHasFibreCapacity)
		}

		ctrl.Finish()
	}
}
