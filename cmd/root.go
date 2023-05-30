// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/config"
)

var (
	// Set via linker flags
	version = "0.0.0"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"

	// Used for flags.
	cfgFile string

	cfg *config.Config = nil

	rootCmd = &cobra.Command{
		Use:   "vand",
		Short: "A Van automation daemon",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Printf("VANd %s, commit %s, built at %s by %s", version, commit, date, builtBy)
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
}

func initConfig() {
	var err error
	if cfg, err = config.NewConfig(cfgFile); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Print(err)
	}
}
