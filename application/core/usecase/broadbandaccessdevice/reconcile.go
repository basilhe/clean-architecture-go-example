package broadbandaccessdevice


const MAX_SERIAL_NUMBER_LENGTH = 25

type OnSuccess interface{
	AuditSuccess()
}

type OnFailure interface{
	AuditFailure()
}

type IGetSerialNumberFromReality interface {
	GetSerialNumber(hostname string) (serialNumber string)
}

type IGetSerialNumberFromModel interface {
	GetSerialNumber(hostname string) (serialNumber string)
}

type IUpdateSerialNumberInModel interface {
	UpdateSerialNumber(hostname, serialNumber string) error
}

type IGetAllDevicesHostname interface {
	GetAllDevicesHostnames() (hostnames []string)
}

type ReconcileBroadbandAccessDeviceUseCase struct {
	IGetAllDevicesHostname
	IGetSerialNumberFromModel
	IGetSerialNumberFromReality
	IUpdateSerialNumberInModel
}

func NewConcileBroadbandAccessDeviceUseCase(iGetAllDevicesHostname IGetAllDevicesHostname, iGetSerialNumberFromModel IGetSerialNumberFromModel, iGetSerialNumberFromReality IGetSerialNumberFromReality, iUpdateSerialNumberInModel IUpdateSerialNumberInModel) *ReconcileBroadbandAccessDeviceUseCase {
	return &ReconcileBroadbandAccessDeviceUseCase{
		iGetAllDevicesHostname,
		iGetSerialNumberFromModel,
		iGetSerialNumberFromReality,
		iUpdateSerialNumberInModel,
	}
}

func (uc *ReconcileBroadbandAccessDeviceUseCase) reconcile(onSuccess OnSuccess, onFailure OnFailure) {
	hostnames := uc.GetAllDevicesHostnames()

	for _, hostname := range hostnames {
		serialNumberInModel := uc.IGetSerialNumberFromModel.GetSerialNumber(hostname)
		serialNumberReality := uc.IGetSerialNumberFromReality.GetSerialNumber(hostname)

		if noSerialNumberInReality(serialNumberReality) || isInvalid(serialNumberReality) {
			onFailure.AuditFailure()
		} else if noSerialNumberInModel(serialNumberInModel) || serialNumberIsDifferentInReality(serialNumberReality, serialNumberInModel) {
			uc.UpdateSerialNumber(hostname, serialNumberReality)
			onSuccess.AuditSuccess()
		}
	}
}

func noSerialNumberInModel(serialNumberInModel string) bool {
	return serialNumberInModel == ""
}

func serialNumberIsDifferentInReality(serialNumberReality string, serialNumberInModel string) bool {
	return serialNumberReality != serialNumberInModel
}

func isInvalid(serialNumberReality string) bool {
	return len(serialNumberReality) > MAX_SERIAL_NUMBER_LENGTH
}

func noSerialNumberInReality(serialNumberReality string) bool {
	return serialNumberReality == ""
}

var reconcileBroadbandAccessDeviceUseCase *ReconcileBroadbandAccessDeviceUseCase
