package controllers

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"time"

	"codicus.ru/deepcool/devices"
	"codicus.ru/deepcool/metrics"
	"github.com/karalabe/hid"
)

type DeviceController struct {
	device      *hid.Device
	metrics     *metrics.Metrics
	deviceInfo  devices.Device
	cpuTdpWatts int
}

func (dc *DeviceController) Initialize() error {
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

	if _, err := dc.device.Write(init); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)

	if _, err := dc.device.Write(dc.Configure(init, false)); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)

	log.Println("Initialization done.")
	return nil
}

func (dc *DeviceController) Configure(initData []byte, hazControlZeroes bool) []byte {
	if hazControlZeroes {
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

func (dc *DeviceController) checksum(data []byte) byte {
	sum := 0
	for _, v := range data[1:16] {
		sum += int(v)
	}
	return byte(sum % 256)
}

func (dc *DeviceController) createStatusPacket(fahrenheit bool) []byte {
	temp := dc.metrics.GetCpuTemp()
	if fahrenheit {
		temp = temp*9/5 + 32
	}

	tempBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tempBytes, math.Float32bits(temp))

	cpuUsage := dc.metrics.GetCpuUsage()
	cpuPower := dc.metrics.GetCpuPower()
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

func (dc *DeviceController) sendPacket(data []byte) error {
	if _, err := dc.device.Write(data); err != nil {
		log.Printf("Write error: %v", err)
		return err
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}

func (dc *DeviceController) SendStatusLoop(sensorName string) {
	for {
		dc.metrics.Update(sensorName, dc.cpuTdpWatts)
		data := dc.createStatusPacket(false)
		if err := dc.sendPacket(data); err != nil {
			break
		}
	}
}

func (dc *DeviceController) Close() {
	dc.device.Close()
}

func NewDeviceController(vid, pid uint16, m *metrics.Metrics, devInfo devices.Device, cTdpWatts int) (*DeviceController, error) {
	devices := hid.Enumerate(vid, pid)
	if len(devices) == 0 {
		return nil, fmt.Errorf("device not found (VID=%04x, PID=%04x)", vid, pid)
	}

	dev, err := devices[0].Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open device: %v", err)
	}

	return &DeviceController{
		device:      dev,
		metrics:     m,
		deviceInfo:  devInfo,
		cpuTdpWatts: cTdpWatts,
	}, nil
}
