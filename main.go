package main

import (
	"log"

	"codicus.ru/deepcool/controllers"
	"codicus.ru/deepcool/devices"
	"codicus.ru/deepcool/metrics"
)

var Version = "1.2.0"
var dcDevice = devices.NewDcLdS360()
var m = &metrics.Metrics{}

func main() {
	log.Println("Version:", Version)
	controller, err := controllers.NewDeviceController(dcDevice.GetVid(), dcDevice.GetPid(), m, dcDevice, 170)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer controller.Close()

	log.Println("Device connected!")

	if err := controller.Initialize(); err != nil {
		log.Fatalf("Initialization error: %v", err)
	}

	controller.SendStatusLoop("k10temp_tctl")
}
