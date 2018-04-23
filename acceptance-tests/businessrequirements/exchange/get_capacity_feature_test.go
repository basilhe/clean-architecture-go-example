package exchange

import (
	"github.com/DATA-DOG/godog"
	"github.com/golang/mock/gomock"
	"github.com/basilhe/tdd/application/core/entity"
	"fmt"
	"strconv"
	"github.com/basilhe/tdd/application/core/usecase/mock"
	usecase "github.com/basilhe/tdd/application/core/usecase/exchange"
)

const NO_PORTS = 0
const EXCHANGE_CODE = "exchange_code"

var mockIGetAvailablePortsOfAllDevicesInExchange *mock.MockIGetAvailablePortsOfAllDevicesInExchange
var mockIDoesExchangeExist *mock.MockIDoesExchangeExist
var useCaseGetCapacityForExchange usecase.IGetCapacityForExchangeUseCase

func parseDeviceType(deviceType string) entity.DeviceType {
	if deviceType == "Fibre" {
		return entity.FIBRE
	} else {
		return entity.ADSL
	}
}
func exchangeExist(exist string) error {
	givenExchange(ctrl, exist == "true")
	return nil
}

func hasPorts(deviceType string, ports int) error {
	devices[deviceIndex] = device(parseDeviceType(deviceType), ports)
	deviceIndex++
	return nil
}

var capacity *entity.Capacity
var err error

func getCapacity() error {
	givenDevices(ctrl, devices...)
	setupUseCase()

	capacity, err = useCaseGetCapacityForExchange.GetCapacity(EXCHANGE_CODE)
	return nil
}

func returnError(hasError string) error {
	if hasError == "true" {
		expectError = true
	}
	if expectError && err == nil {
		return fmt.Errorf("Error should occur")
	} else if (hasError == "false" && err != nil) {
		return fmt.Errorf("Error should not occur")
	}
	return nil
}

func hasCapacity(deviceType, hasCapacity string) error {
	if expectError {
		return nil
	}
	if (deviceType == "Fibre") {
		if (strconv.FormatBool(capacity.HasFibreCapacity()) != hasCapacity) {
			return fmt.Errorf("HasFibreCapacity doesn't match")
		}
	} else if (deviceType == "Adsl") {
		if (strconv.FormatBool(capacity.HasAdslCapacity()) != hasCapacity) {
			return fmt.Errorf("HasAdslCapacity doesn't match")
		}
	}
	return nil
}

var ctrl *gomock.Controller
var devices []*entity.BroadbandAccessDevice
var deviceIndex int
var expectError bool

func FeatureContext(s *godog.Suite) {
	s.Step(`^exchange exist "([^"]*)"$`, exchangeExist)
	s.Step(`^has "([^"]*)" "(\d+)" ports$`, hasPorts)
	s.Step(`^get capacity$`, getCapacity)
	s.Step(`^return error "([^"]*)"$`, returnError)
	s.Step(`^has "([^"]*)" capacity "([^"]*)"$`, hasCapacity)

	s.BeforeScenario(func(interface{}) {
		ctrl = gomock.NewController(&cancelReporter{})
		devices = make([]*entity.BroadbandAccessDevice, 2)
		deviceIndex = 0
		expectError = false
	})

	s.AfterScenario(func(i interface{}, e error) {
		ctrl.Finish()
	})
}

type cancelReporter struct {
	cancel func()
}

func (r *cancelReporter) Errorf(format string, args ...interface{}) {
	fmt.Errorf(format, args)
}
func (r *cancelReporter) Fatalf(format string, args ...interface{}) {
	fmt.Errorf(format, args)
	//defer r.cancel()
}


func setupUseCase() {
	useCaseGetCapacityForExchange = usecase.NewGetCapacityForExchangeUseCase(mockIGetAvailablePortsOfAllDevicesInExchange, mockIDoesExchangeExist)
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