package devices

import (
	"math"
	"time"

	"codicus.ru/deepcool/metrics"
	"encoding/binary"
)

type DcLdS360 struct {
	vid uint16 // Vendor ID
	pid uint16 // Product ID
}

func (dc *DcLdS360) GetVid() uint16 {
	return dc.vid
}

func (dc *DcLdS360) GetPid() uint16 {
	return dc.pid
}

func (dc *DcLdS360) GetStatusPacket() []byte {
	init := make([]byte, 64)
	init[0] = 0x10
	init[1] = 0x68
	init[2] = 0x01
	init[3] = 0x01
	init[4] = 0x02
	init[5] = 0x03
	init[6] = 0x01
	init[7] = 0x70
	init[8] = 0x16

	return init
}

func (dc *DcLdS360) GetConfigurePacket(initData []byte, hazControlZeros bool) []byte {
	if hazControlZeros {
		initData[5] = 0x02
		initData[7] = 0x6F
		time.Sleep(50 * time.Millisecond)
	} else {
		initData[5] = 0x02
		initData[6] = 0x00
		initData[7] = 0x6E
		time.Sleep(50 * time.Millisecond)
	}
	return initData
}

func (dc *DcLdS360) checksum(data []byte) byte {
	sum := 0
	for _, v := range data[1:16] {
		sum += int(v)
	}
	return byte(sum % 256)
}

func (dc *DcLdS360) GetDataPacket(fahrenheit bool, m metrics.Metrics) []byte {
	temp := m.GetCpuTemp()
	if fahrenheit {
		temp = temp*9/5 + 32
	}

	tempBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tempBytes, math.Float32bits(temp))

	cpuUsage := m.GetCpuUsage()
	cpuPower := m.GetCpuPower()
	powerBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(powerBytes, uint16(cpuPower))

	statusData := make([]byte, 64)
	statusData[0] = 0x10
	statusData[1] = 0x68
	statusData[2] = 0x01
	statusData[3] = 0x01
	statusData[4] = 0x0B
	statusData[5] = 0x01
	statusData[6] = 0x02
	statusData[7] = 0x05
	statusData[8] = powerBytes[0]
	statusData[9] = powerBytes[1]
	if fahrenheit {
		statusData[10] = 1
	} else {
		statusData[10] = 0
	}
	copy(statusData[11:15], tempBytes)
	statusData[15] = byte(cpuUsage)
	statusData[16] = dc.checksum(statusData)
	statusData[17] = 0x16

	return statusData
}

func NewDcLdS360() *DcLdS360 {
	return &DcLdS360{
		vid: 0x3633,
		pid: 0x000A,
	}
}
