package app

import (
	"log"

	"github.com/codicuz/deepcool/v2/controllers"
	"github.com/codicuz/deepcool/v2/devices"
	"github.com/codicuz/deepcool/v2/metrics"
)

func Run(sensorName string, deviceModel string, output bool, cpuTdpWatts uint16, interval uint16) {
	var dcDevice devices.Device
	switch deviceModel {
	case "dc_ld_s360":
		dcDevice = devices.NewDcLdS360()
	default:
		log.Fatalf("Unknown device: %s", deviceModel)
	}

	devices.NewDcLdS360()
	var m = &metrics.Metrics{}

	controller, err := controllers.NewDeviceController(dcDevice.GetVid(), dcDevice.GetPid(), m, dcDevice, cpuTdpWatts, interval)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer controller.Close()

	log.Println("Device connected!")

	if err := controller.Initialize(false); err != nil {
		log.Fatalf("Initialization error: %v", err)
	}

	// controller.SendStatusLoop("k10temp_tctl")
	controller.SendStatusLoop(sensorName, output)
}
