package main

import (
	// "gioui.org/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(displayCmd)
}

var displayCmd = &cobra.Command{
	Use:   "display",
	Short: "Start the display agent",
	Run:   runDisplay,
}

func runDisplay(cmd *cobra.Command, args []string) {

	// dev, err := mockup.New()
	// if err != nil {
	// 	log.Fatalf("Failed to open display: %s", err)
	// }

	// go func() {
	// 	showCanvas(dev)
	// 	time.Sleep(1 * time.Second)

	// 	playGif(dev)
	// }()

	// app.Main()
}
