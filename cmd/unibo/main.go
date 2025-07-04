// SPDX-FileCopyrightText: 2024 - 2025 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: MIT

package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "unibo",
	Short: "A CLI to interact with the University of Bologna",

	RunE: rootRun,
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}

func rootRun(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
