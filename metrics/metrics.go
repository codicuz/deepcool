package metrics

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
)

type Metrics struct {
	cpuTemp  float32
	cpuUsage int
	cpuPower int
}

func (m *Metrics) getCpuTemp(sensorName string) float32 {
	sensors, err := host.SensorsTemperatures()
	if err != nil {
		return 0.0
	}
	for _, s := range sensors {
		if s.SensorKey == sensorName {
			return float32(s.Temperature)
		}
	}
	return 0.0
}

func (m *Metrics) getCpuUsage() int {
	percentages, err := cpu.Percent(time.Second/10, false)
	if err != nil || len(percentages) == 0 {
		return 0
	}
	return int(percentages[0])
}

func (m *Metrics) getCpuPower(usagePercent int, cpuTdpWatts int) int {
	return int(cpuTdpWatts) * usagePercent / 100
}

func (m *Metrics) Update(sensorName string, cpuTdpWatts int) {
	m.cpuTemp = m.getCpuTemp(sensorName)
	m.cpuUsage = m.getCpuUsage()
	m.cpuPower = m.getCpuPower(m.cpuUsage, cpuTdpWatts)
}

func (m *Metrics) GetCpuTemp() float32 {
	return m.cpuTemp
}

func (m *Metrics) GetCpuUsage() int {
	return m.cpuUsage
}

func (m *Metrics) GetCpuPower() int {
	return m.cpuPower
}
