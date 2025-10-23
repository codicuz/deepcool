package controllers

import (
	"fmt"
	"log"
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

func (dc *DeviceController) Initialize(hasControlZeros bool) error {
	init := dc.deviceInfo.GetStatusPacket()

	if _, err := dc.device.Write(init); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)

	if _, err := dc.device.Write(dc.Configure(init, hasControlZeros)); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)

	log.Println("Initialization done.")
	return nil
}

func (dc *DeviceController) Configure(initData []byte, hazControlZeros bool) []byte {
	iData := dc.deviceInfo.GetConfigurePacket(initData, hazControlZeros)

	return iData
}

func (dc *DeviceController) createStatusPacket(fahrenheit bool) []byte {
	statusData := dc.deviceInfo.GetDataPacket(fahrenheit, *dc.metrics)

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

func (dc *DeviceController) SendStatusLoop(sensorName string, output bool) {
	for {
		dc.metrics.Update(sensorName, dc.cpuTdpWatts)
		data := dc.createStatusPacket(false)
		if err := dc.sendPacket(data); err != nil {
			break
		}
		if output {
			log.Printf("Sending status: Temp=%.2f°C, Usage=%d%%, Power=%dW", dc.metrics.GetCpuTemp(), dc.metrics.GetCpuUsage(), dc.metrics.GetCpuPower())
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
