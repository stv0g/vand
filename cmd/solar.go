package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/devices/solar/renogy"
	"github.com/stv0g/vand/pkg/pb"
	"github.com/stv0g/vand/pkg/types"
)

func init() {
	rootCmd.AddCommand(solarCmd)
}

var solarCmd = &cobra.Command{
	Use:   "solar",
	Short: "Start the photo-voltaic agent",
	Run:   runSolar,
}

func runSolar(cmd *cobra.Command, args []string) {
	d, err := renogy.NewDevice(cfg.Solar.Address)

	sts, err := d.GetState()
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(sts)

	cfg, err := d.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(cfg)

	sum := pb.StateUpdatePoint{
		Solar: sts,
	}

	f := types.Flatten(sum)
	for k, v := range f {
		fmt.Printf("%s = %v\n", k, v)

	}
}
