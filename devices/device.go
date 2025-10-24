package devices

import "github.com/codicuz/deepcool/metrics"

type Device interface {
	GetVid() uint16
	GetPid() uint16
	GetStatusPacket() []byte
	GetConfigurePacket([]byte, bool) []byte
	checksum([]byte) byte
	GetDataPacket(fahrenheit bool, m metrics.Metrics) []byte
}
