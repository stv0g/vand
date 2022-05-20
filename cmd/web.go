package main

import (
	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/web"
)

func init() {
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server",
	Run:   runWeb,
}

func runWeb(cmd *cobra.Command, args []string) {
	web.Run(cfg, version, commit, date)
}
