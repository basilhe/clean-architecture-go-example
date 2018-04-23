package exchange

import (
	"github.com/basilhe/tdd/application/core/entity"
	"errors"
)

var ErrorExchangeNotFound = errors.New("Exchange not found")

type IGetAvailablePortsOfAllDevicesInExchange interface {
	GetAvailablePortsOfAllDevicesInExchange(exchange string) ([]*entity.BroadbandAccessDevice, error)
}

type IDoesExchangeExist interface {
	DoesExchangeExist(exchange string) bool
}

type IGetCapacityForExchangeUseCase interface {
	GetCapacity(exchange string) (*entity.Capacity, error)
}

type getCapacityForExchangeUseCase struct {
	IGetAvailablePortsOfAllDevicesInExchange
	IDoesExchangeExist
}

func NewGetCapacityForExchangeUseCase(iGetAvailablePortsOfAllDevicesInExchange IGetAvailablePortsOfAllDevicesInExchange, iDoesExchangeExist IDoesExchangeExist) IGetCapacityForExchangeUseCase {
	return &getCapacityForExchangeUseCase{iGetAvailablePortsOfAllDevicesInExchange, iDoesExchangeExist}

}

func (useCase *getCapacityForExchangeUseCase) GetCapacity(exchange string) (*entity.Capacity, error){
	exchangeExist := useCase.DoesExchangeExist(exchange)
	if !exchangeExist {
		return nil, errors.New("Exchange not found")
	}
	devices, _ := useCase.GetAvailablePortsOfAllDevicesInExchange(exchange)

	hasAdslCapacity := hasCapacityFor(devices, entity.ADSL)
	hasFibreCapacity := hasCapacityFor(devices, entity.FIBRE)
	return &entity.Capacity{hasAdslCapacity, hasFibreCapacity}, nil
}

const MINIMUM_NUMBER_OF_PORTS = 5

func hasCapacityFor(devices []*entity.BroadbandAccessDevice, deviceType entity.DeviceType) bool {
	availablePorts := countAvailablePortsForType(devices, deviceType)

	return availablePorts >= MINIMUM_NUMBER_OF_PORTS;
}

func countAvailablePortsForType(devices []*entity.BroadbandAccessDevice, deviceType entity.DeviceType) int {
	var availablePorts int
	for _, d := range devices {
		if d.DeviceType == deviceType {
			availablePorts += d.AvailablePorts
		}
	}
	return availablePorts
}
