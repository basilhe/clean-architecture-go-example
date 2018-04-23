package broadbandaccessdevice

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/golang/mock/gomock"
	"github.com/basilhe/tdd/application/core/entity"
	"github.com/basilhe/tdd/application/core/usecase/mock"
)


var getDetailsUseCase *GetDetailsUseCase
var mockIGetDeviceDetails *mock.MockIGetDeviceDetails

func TestReturnsDeviceDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var expectedDevice *entity.BroadbandAccessDevice
	expectedDevice = givenADeviceIsFound(ctrl)
	setupGetDetailsUseCase()

	actualDevice, err := getDetailsUseCase.GetDeviceDetails("hostname1")

	assert.NoError(t, err)
	assert.EqualValues(t, expectedDevice, actualDevice)
}

func TestErrorWhenDeviceIsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenADeviceIsNotFound(ctrl)
	setupGetDetailsUseCase()

	var aDevice *entity.BroadbandAccessDevice
	aDevice, err := getDetailsUseCase.GetDeviceDetails("hostname1")

	assert.Nil(t, aDevice, "Device should not be found")
	assert.EqualError(t, err, "Device is not found")
}

func setupGetDetailsUseCase() {
	getDetailsUseCase = NewGetDetailsUseCase(mockIGetDeviceDetails)
}

func givenADeviceIsFound(ctrl *gomock.Controller) *entity.BroadbandAccessDevice {
	expectedDevice := &entity.BroadbandAccessDevice{"hostname1", "serialNumber", entity.ADSL, 0}
	mockIGetDeviceDetails = mock.NewMockIGetDeviceDetails(ctrl)
	mockIGetDeviceDetails.EXPECT().GetDetails("hostname1").Return(expectedDevice, nil)

	return expectedDevice
}

func givenADeviceIsNotFound(ctrl *gomock.Controller) {
	mockIGetDeviceDetails = mock.NewMockIGetDeviceDetails(ctrl)
	mockIGetDeviceDetails.EXPECT().GetDetails("hostname1").Return(nil, ErrorDeviceIsNotFound)
}
