package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "unibo",
	Short: "A CLI to interact with the University of Bologna",

	Run: rootRun,
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}

func rootRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}
