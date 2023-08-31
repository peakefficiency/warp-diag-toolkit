/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/peakefficiency/warp-diag-toolkit/cli"
	"github.com/peakefficiency/warp-diag-toolkit/diag"
	"github.com/peakefficiency/warp-diag-toolkit/info"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		diag.ZipPath = args[0]
		contents, err := diag.ExtractToMemory(diag.ZipPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		info.GetInfo(diag.ZipPath, contents)

		if cli.Debug {
			fmt.Println("Files in zip:")
			fmt.Println()
			for filename := range contents {
				fmt.Println(filename)
			}
			if content, ok := contents["connectivity.txt"]; ok {
				fmt.Println()
				fmt.Println("Debug testing connectivity.txt:")
				fmt.Println(string(content.Data))
			}
		}
		if cli.Debug {
			fmt.Println("Debug check info read: ")
			fmt.Printf("debug Platform type: %s\n", info.Info.PlatformType)
			fmt.Printf("debug Split tunnel mode: %s\n", info.Info.SplitTunnelMode)
			fmt.Printf("debug Split tunnel list: \n%s", info.Info.SplitTunnelList)
			fmt.Printf("debug Fallback domains: \n%s", info.Info.FallbackDomains)

		}
		// Print Markdown output
		fmt.Println("info called")
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
