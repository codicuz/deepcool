package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"time"

	"codicus.ru/deepcool/devices"
	"github.com/karalabe/hid"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
)

var Version = "1.1.0"

var dc = devices.NewDcLd360()

// DeviceController — структура для работы с устройством
type DeviceController struct {
	device *hid.Device
}

// Создание нового контроллера
func NewDeviceController(vid, pid uint16) (*DeviceController, error) {
	devices := hid.Enumerate(vid, pid)
	if len(devices) == 0 {
		return nil, fmt.Errorf("device not found (VID=%04x, PID=%04x)", vid, pid)
	}

	dev, err := devices[0].Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open device: %v", err)
	}

	return &DeviceController{device: dev}, nil
}

// Закрытие устройства
func (dc *DeviceController) Close() {
	dc.device.Close()
}

// Вычисление контрольной суммы
func checksum(data []byte) byte {
	sum := 0
	for _, v := range data[1:16] {
		sum += int(v)
	}
	return byte(sum % 256)
}

// Получение температуры CPU
func getCPUTemp() float32 {
	sensors, err := host.SensorsTemperatures()
	if err != nil {
		return 0.0
	}
	for _, s := range sensors {
		if s.SensorKey == "k10temp_tctl" {
			return float32(s.Temperature)
		}
	}
	return 0.0
}

// Получение загрузки CPU
func getCPUUsage() int {
	percentages, err := cpu.Percent(time.Second/10, false)
	if err != nil || len(percentages) == 0 {
		return 0
	}
	return int(percentages[0])
}

// Расчет потребляемой мощности CPU
func getCPUPower(usagePercent int) int {
	return int(dc.GetCpuTdpWatts()) * usagePercent / 100
}

// Создание пакета статуса
func (dc *DeviceController) createStatusPacket(fahrenheit bool) []byte {
	temp := getCPUTemp()
	if fahrenheit {
		temp = temp*9/5 + 32
	}

	tempBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tempBytes, math.Float32bits(temp))

	cpuUsage := getCPUUsage()
	cpuPower := getCPUPower(cpuUsage)
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
	statusData[16] = checksum(statusData)
	statusData[17] = 0x16

	return statusData
}

// Инициализация устройства
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
	}
	return initData
}

func (dc *DeviceController) sendPacket(data []byte) error {
	if _, err := dc.device.Write(data); err != nil {
		log.Printf("Write error: %v", err)
		return err
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}

// Отправка пакета статуса
func (dc *DeviceController) SendStatusLoop() {
	for {
		data := dc.createStatusPacket(false)
		if err := dc.sendPacket(data); err != nil {
			break
		}
	}
}

func main() {
	log.Println("Version:", Version)
	controller, err := NewDeviceController(dc.GetVid(), dc.GetPid())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer controller.Close()

	log.Println("Device connected!")

	if err := controller.Initialize(); err != nil {
		log.Fatalf("Initialization error: %v", err)
	}

	controller.SendStatusLoop()
}
