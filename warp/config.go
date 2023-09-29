package warp

import (
	_ "embed"
	"errors"
	"fmt"
	"io"

	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var WdcConf WDCConfig
var SaveReport, Verbose, Debug, Offline, Plain bool

type WDCConfig struct {
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

func GetOrLoadConfigWdc() {

	if Offline {

		LocalConfig()
		LoadConfig()
		return
	}

	RemoteConfig("https://warp-diag-checker.pages.dev/wdc-config.yaml")
	LoadConfig()

}

func RemoteConfig(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(errors.New("unable to get remote config"))
	}

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
	var config WDCConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println("Failed to parse YAML file:", err)

	}
	WdcConf = config

}
