/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

// saveconfigCmd represents the saveconfig command
var saveconfigCmd = &cobra.Command{
	Use:   "saveconfig",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		var yamlFile []byte

		resp, err := http.Get("https://warp-diag-checker.pages.dev/wdc-config.yaml")
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to download YAML file: HTTP %d\n", resp.StatusCode)
		}
		// read the response body
		yamlFile, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read response body:", err)

			return
		}
		// save the YAML file to the user's home folder
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Failed to get current user:", err)
			return
		}
		configPath := filepath.Join(usr.HomeDir, "wdc-config.yaml")
		err = os.WriteFile(configPath, yamlFile, 0644)
		if err != nil {
			fmt.Println("Failed to save YAML file:", err)
			return
		}
		fmt.Println("Configuration saved to:", configPath)

	},
}

func init() {
	rootCmd.AddCommand(saveconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// saveconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// saveconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
