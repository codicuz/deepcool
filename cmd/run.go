package cmd

import (
	"github.com/codicuz/deepcool/v2/app"
	"github.com/spf13/cobra"
)

var tempSensorName string
var deviceModel string
var cpuTdpWatts uint16
var interval uint16
var output bool

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run application",
	Long:  "Run DeepCool application",
	Run: func(cmd *cobra.Command, args []string) {
		app.Run(tempSensorName, deviceModel, output, cpuTdpWatts, interval)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&tempSensorName, "temperature-sensor", "t", "", "Temperature sensor name [k10temp_tctl]")
	runCmd.Flags().StringVarP(&deviceModel, "device-model", "d", "", "Device model [dc_ld_s360]")
	
	runCmd.Flags().Uint16VarP(&cpuTdpWatts, "cpu-tdp-watts", "c", 170, "Thermal Design Power (Watts)")
	runCmd.Flags().Uint16VarP(&interval, "interval", "i", 500, "Interval between packet sends")
	runCmd.Flags().BoolVarP(&output, "output", "o", false, "Enable output")

	runCmd.MarkFlagRequired("temperature-sensor")
	runCmd.MarkFlagRequired("device-model")
}
