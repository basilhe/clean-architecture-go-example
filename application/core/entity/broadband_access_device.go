package entity

type BroadbandAccessDevice struct {
	Hostname       string
	SerialNumber   string
	DeviceType     DeviceType
	AvailablePorts int
}
