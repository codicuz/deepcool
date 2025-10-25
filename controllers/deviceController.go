package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/codicuz/deepcool/v2/devices"
	"github.com/codicuz/deepcool/v2/metrics"
	"github.com/karalabe/hid"
)

type DeviceController struct {
	device      *hid.Device
	metrics     *metrics.Metrics
	deviceInfo  devices.Device
	timeout     time.Duration
	cpuTdpWatts uint16
}

func (dc *DeviceController) Initialize(hasControlZeros bool) error {
	init := dc.deviceInfo.GetStatusPacket()

	if _, err := dc.device.Write(init); err != nil {
		return err
	}
	time.Sleep(dc.timeout)

	if _, err := dc.device.Write(dc.Configure(init, hasControlZeros)); err != nil {
		return err
	}
	time.Sleep(dc.timeout)

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
	time.Sleep(dc.timeout)
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

func NewDeviceController(vid, pid uint16, m *metrics.Metrics, devInfo devices.Device, cTdpWatts uint16, tOut uint16) (*DeviceController, error) {
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
		timeout: time.Duration(tOut) * time.Millisecond,
	}, nil
}
