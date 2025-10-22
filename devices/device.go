package devices

type Device interface {
	GetVid() uint16
	GetPid() uint16
}