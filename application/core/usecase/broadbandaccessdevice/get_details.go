package broadbandaccessdevice

import (
	"github.com/basilhe/tdd/application/core/entity"
	"errors"
)

type IGetDeviceDetails interface {
	GetDetails(hostname string) (*entity.BroadbandAccessDevice, error)
}

var ErrorDeviceIsNotFound = errors.New("Device is not found")

type GetDetailsUseCase struct {
	IGetDeviceDetails
}

func NewGetDetailsUseCase(iGetDeviceDetails IGetDeviceDetails) *GetDetailsUseCase {
	return &GetDetailsUseCase{iGetDeviceDetails}
}

func (uc *GetDetailsUseCase) GetDeviceDetails(hostname string) (aDevice *entity.BroadbandAccessDevice, err error) {
	return uc.GetDetails(hostname)
}
