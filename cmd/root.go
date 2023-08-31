/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/peakefficiency/warp-diag-toolkit/diag"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wdc",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wdc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().BoolVarP(&diag.SaveReport, "report", "r", false, "Save the generated report in the local folder")
	rootCmd.PersistentFlags().BoolVarP(&diag.Plain, "plain", "p", false, "Output the report in plain markdown")
	rootCmd.PersistentFlags().BoolVarP(&diag.Verbose, "verbose", "v", false, "Increase output verbosity")
	rootCmd.PersistentFlags().BoolVarP(&diag.Debug, "debug", "", false, "Enable debug mode (not implemented yet)")
	rootCmd.PersistentFlags().BoolVarP(&diag.Offline, "offline", "o", false, "Force the use of the local YAML cache file")
}
