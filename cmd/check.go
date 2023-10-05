/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/peakefficiency/warp-diag-toolkit/warp"
	"github.com/spf13/cobra"
)

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
		warp.ZipPath = args[0]
		contents, err := warp.ExtractToMemory(warp.ZipPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		warp.GetOrLoadConfig(warp.WdcConfig)

		info := contents.GetInfo(warp.ZipPath)

		warp.NewPrinter().PrintCheckResult(info.VersionCheck())

		warp.NewPrinter().PrintCheckResult(info.SplitTunnelCheck())

		warp.NewPrinter().PrintCheckResult(info.DefaultExcludeCheck())

		contents.LogSearch(info)

		warp.NewPrinter().PrintString(warp.ReportLogSearch(warp.LogSearchOutput))

		warp.NewPrinter().PrintString(info.ReportInfo())

		fmt.Println("check called")
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	rootCmd.PersistentFlags().BoolVarP(&warp.SaveReport, "report", "r", false, "Save the generated report in the local folder")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
