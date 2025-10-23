package app

import (
	"log"

	"codicus.ru/deepcool/controllers"
	"codicus.ru/deepcool/devices"
	"codicus.ru/deepcool/metrics"
)

func Run(sensorName string, deviceModel string, output bool) {
	var dcDevice devices.Device
	switch deviceModel {
	case "dc_ld_s360":
		dcDevice = devices.NewDcLdS360()
	default:
		log.Fatalf("Unknown device: %s", deviceModel)
	}

	devices.NewDcLdS360()
	var m = &metrics.Metrics{}

	controller, err := controllers.NewDeviceController(dcDevice.GetVid(), dcDevice.GetPid(), m, dcDevice, 170)
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
