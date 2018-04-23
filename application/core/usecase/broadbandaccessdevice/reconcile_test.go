package broadbandaccessdevice

import (
	"github.com/basilhe/tdd/application/core/usecase/mock"
	"github.com/golang/mock/gomock"
	"testing"
)

var mockOnSuccess *mock.MockOnSuccess
var mockOnFailure *mock.MockOnFailure
var mockIGetSerialNumberFromReality *mock.MockIGetSerialNumberFromReality
var mockIGetSerialNumberFromModel *mock.MockIGetSerialNumberFromModel
var mockIGetAllDevicesHostname *mock.MockIGetAllDevicesHostname
var mockIUpdateSerialNumberInModel *mock.MockIUpdateSerialNumberInModel

func TestNothingToUpdateWhenModelAndRealityAreTheSame(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenThereIsADeviceWithHostname(ctrl, "hostname1")
	givenDeviceHasSerialNumberInModel(ctrl, "hostname1", "serialNumber1")
	givenDeviceHasSerialNumberInReality(ctrl, "hostname1", "serialNumber1")
	expectedNothingHasBeenUpdated(ctrl)
	setupConcileUseCase(ctrl)

	reconcileBroadbandAccessDeviceUseCase.reconcile(mockOnSuccess, mockOnFailure)
}

func TestUpdateSerialNumberWhenRealityIsDifferentFromModel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenThereIsADeviceWithHostname(ctrl, "hostname1")
	givenDeviceHasSerialNumberInModel(ctrl, "hostname1", "serialNumber1")
	givenDeviceHasSerialNumberInReality(ctrl, "hostname1", "newSerialNumber")

	expectTheDeviceHasBeenUpdated(ctrl, "hostname1", "newSerialNumber")
	setupConcileUseCase(ctrl)

	expectASuccessHasBeenAudited(ctrl)

	reconcileBroadbandAccessDeviceUseCase.reconcile(mockOnSuccess, mockOnFailure)
}

func TestUpdateSerialNumberWhenDeviceDoesNotHaveOneInModel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenThereIsADeviceWithHostname(ctrl, "hostname1")
	givenDeviceHasSerialNumberInModel(ctrl, "hostname1", "")
	givenDeviceHasSerialNumberInReality(ctrl, "hostname1", "newSerialNumber")

	expectTheDeviceHasBeenUpdated(ctrl, "hostname1", "newSerialNumber")
	setupConcileUseCase(ctrl)

	expectASuccessHasBeenAudited(ctrl)

	reconcileBroadbandAccessDeviceUseCase.reconcile(mockOnSuccess, mockOnFailure)
}

func TestAuditsSuccessWhenUpdatesTheModel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenThereIsADeviceWithHostname(ctrl, "hostname1")
	givenDeviceHasSerialNumberInModel(ctrl, "hostname1", "serialNumber1")
	givenDeviceHasSerialNumberInReality(ctrl, "hostname1", "newSerialNumber")

	expectTheDeviceHasBeenUpdated(ctrl, "hostname1", "newSerialNumber")
	setupConcileUseCase(ctrl)

	expectASuccessHasBeenAudited(ctrl)

	reconcileBroadbandAccessDeviceUseCase.reconcile(mockOnSuccess, mockOnFailure)
}

func TestAuditsFailureWhenSerialNumberFromRealityIsLongerThanUsual(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenThereIsADeviceWithHostname(ctrl, "hostname1")
	givenDeviceHasSerialNumberInModel(ctrl, "hostname1", "serialNumber1")
	givenDeviceHasSerialNumberInReality(ctrl, "hostname1", "longerThanAllowedSerialNumber")

	expectedNothingHasBeenUpdated(ctrl)
	setupConcileUseCase(ctrl)

	expectAFailureHasBeenAudited(ctrl)

	reconcileBroadbandAccessDeviceUseCase.reconcile(mockOnSuccess, mockOnFailure)
}

func TestAuditsFailureWhenItCantReconcileSerialNumberFromReality(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	givenThereIsADeviceWithHostname(ctrl, "hostname1")
	givenDeviceHasSerialNumberInModel(ctrl, "hostname1", "serialNumber1")
	givenThereIsAProblemRetrievingTheSerialNumberFromReality(ctrl, "hostname1")

	expectedNothingHasBeenUpdated(ctrl)
	setupConcileUseCase(ctrl)

	expectAFailureHasBeenAudited(ctrl)

	reconcileBroadbandAccessDeviceUseCase.reconcile(mockOnSuccess, mockOnFailure)
}

func givenDeviceHasSerialNumberInReality(ctrl *gomock.Controller, hostname string, serialNumber string) {
	mockIGetSerialNumberFromReality = mock.NewMockIGetSerialNumberFromReality(ctrl)
	mockIGetSerialNumberFromReality.EXPECT().GetSerialNumber(hostname).Return(serialNumber)
}

func givenThereIsAProblemRetrievingTheSerialNumberFromReality(ctrl *gomock.Controller, hostname string) {
	mockIGetSerialNumberFromReality = mock.NewMockIGetSerialNumberFromReality(ctrl)
	mockIGetSerialNumberFromReality.EXPECT().GetSerialNumber(hostname).Return("")
}

func givenDeviceHasSerialNumberInModel(ctrl *gomock.Controller, hostname string, serialNumber string) {
	mockIGetSerialNumberFromModel = mock.NewMockIGetSerialNumberFromModel(ctrl)
	mockIGetSerialNumberFromModel.EXPECT().GetSerialNumber(hostname).Return(serialNumber)
}

func givenThereIsADeviceWithHostname(ctrl *gomock.Controller, hostname ...string) {
	mockIGetAllDevicesHostname = mock.NewMockIGetAllDevicesHostname(ctrl)
	mockIGetAllDevicesHostname.EXPECT().GetAllDevicesHostnames().Return(hostname)
}

func expectedNothingHasBeenUpdated(ctrl *gomock.Controller) {
	mockIUpdateSerialNumberInModel = mock.NewMockIUpdateSerialNumberInModel(ctrl)
}

func expectTheDeviceHasBeenUpdated(ctrl *gomock.Controller, hostname string, serialNumber string) {
	mockIUpdateSerialNumberInModel = mock.NewMockIUpdateSerialNumberInModel(ctrl)
	mockIUpdateSerialNumberInModel.EXPECT().UpdateSerialNumber(hostname, serialNumber).Return(nil)
}

func expectAFailureHasBeenAudited(ctrl *gomock.Controller) {
	mockOnFailure = mock.NewMockOnFailure(ctrl)
	mockOnFailure.EXPECT().AuditFailure()
}

func expectASuccessHasBeenAudited(ctrl *gomock.Controller) {
	mockOnSuccess = mock.NewMockOnSuccess(ctrl)
	mockOnSuccess.EXPECT().AuditSuccess()
}

func setupConcileUseCase(ctrl *gomock.Controller) {
	mockOnSuccess = mock.NewMockOnSuccess(ctrl)

	reconcileBroadbandAccessDeviceUseCase = NewConcileBroadbandAccessDeviceUseCase(mockIGetAllDevicesHostname, mockIGetSerialNumberFromModel, mockIGetSerialNumberFromReality, mockIUpdateSerialNumberInModel)
}

