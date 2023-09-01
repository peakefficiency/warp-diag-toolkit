/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/peakefficiency/warp-diag-toolkit/checks"
	"github.com/peakefficiency/warp-diag-toolkit/config"
	"github.com/peakefficiency/warp-diag-toolkit/data"
	"github.com/peakefficiency/warp-diag-toolkit/output"
	"github.com/spf13/cobra"
)

var ZipPath string

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data.ZipPath = args[0]
		contents, err := data.ExtractToMemory(data.ZipPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		config.GetAndLoadConfig()

		checks.LogSearch(contents)
		searchreport, err := output.ReportLogSearch(checks.LogSearchOutput)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(searchreport)

		info := data.GetInfo(data.ZipPath, contents)

		inforeport, err := output.ReportInfo(info)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(inforeport)

		fmt.Println("check called")
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
