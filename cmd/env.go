// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package main

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Start the environment agent",
	Run:   runEnv,
}

func runEnv(cmd *cobra.Command, args []string) {
}
