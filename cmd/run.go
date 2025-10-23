package cmd

import (
	"codicus.ru/deepcool/app"
	"github.com/spf13/cobra"
)

var tempSensorName string
var deviceModel string
var output bool

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run application",
	Long: "Run DeepCool application",
	Run: func(cmd *cobra.Command, args []string) {
		app.Run(tempSensorName, deviceModel, output)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&tempSensorName, "temperature-sensor", "t", "", "Temperature sensor name [k10temp_tctl]")
	runCmd.MarkFlagRequired("temperature-sensor")
	runCmd.Flags().StringVarP(&deviceModel, "device-model", "d", "", "Device model [dc_ld_s360]")
	runCmd.MarkFlagRequired("device-model")
	runCmd.Flags().BoolVarP(&output, "output", "o", false, "Enable output")
}
