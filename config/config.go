package config

import (
	_ "embed"
	"fmt"
	"io"

	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/peakefficiency/warp-diag-toolkit/cli"

	"gopkg.in/yaml.v2"
)

var Conf Config

type Config struct {
	AppReleaseVersion  string   `yaml:"wdc_latest_version"`
	ConfigVersion      string   `yaml:"config_version"`
	BadVersions        []string `yaml:"bad_versions"`
	LogPatternsByIssue []struct {
		SearchFile string `yaml:"search_file"`
		Issue      map[string]struct {
			SearchTerms []string `yaml:"search_term"`
		} `yaml:"issue_type"`
		ReplyType map[string]struct {
			Message string `yaml:"message"`
		} `yaml:"reply_type"`
	} `yaml:"log_patterns_by_issue"`
	ReplyByIssueType map[string]struct {
		Message string `yaml:"message"`
	} `yaml:"reply_by_issue_type"`
}

//go:embed wdc-config.yaml
var embeddedConfig []byte
var yamlFile []byte
var err error

func LocalConfig() {
	configPath := "./wdc-config.yaml"
	yamlFile, err = os.ReadFile(configPath)
	if err != nil {
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Failed to get current user:", err)
			return
		}
		configPath := filepath.Join(usr.HomeDir, "wdc-config.yaml")
		yamlFile, err = os.ReadFile(configPath)
		if err != nil {
			fmt.Println("Failed to read local YAML file:", err)
			yamlFile = embeddedConfig // use the embedded config as fallback
		}
	}
}

func GetOrLoadConfig() {

	if cli.Offline {

		LocalConfig()
		LoadConfig()
		return
	}

	RemoteConfig()
	LoadConfig()

	return
}

func RemoteConfig() {
	resp, err := http.Get("https://warp-diag-checker.pages.dev/wdc-config.yaml")

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to download YAML file: HTTP %d\n", resp.StatusCode)

		LocalConfig()
	}
	// read the response body
	yamlFile, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		LocalConfig()
		return
	}

}

func LoadConfig() {
	// Load Config from the YAML
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println("Failed to parse YAML file:", err)

	}
	Conf = config

	if cli.Debug {
		fmt.Println("Config Version", Conf.ConfigVersion)
		// print the log patterns by issue
		for _, logPattern := range Conf.LogPatternsByIssue {
			fmt.Printf("Search File: %s\n", logPattern.SearchFile)
			for issueTypeName, issue := range logPattern.Issue {
				fmt.Printf("Issue Type: %s\n", issueTypeName)
				for _, searchTerm := range issue.SearchTerms {
					fmt.Printf("Search Term: %s\n", searchTerm)
				}
				fmt.Println()
			}
		}

		// print all the replies for all available issue types
		fmt.Println("Replies:")
		for issueTypeName, reply := range Conf.ReplyByIssueType {
			fmt.Printf("Issue Type: %s\n", issueTypeName)
			fmt.Printf("Reply: %s\n", reply.Message)
			fmt.Println()
		}

		// print a diagnostic message
		fmt.Println("Config loaded successfully")
	}
	return
}
