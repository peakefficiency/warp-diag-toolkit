package data

import (
	"archive/zip"
	"bytes"
	"io"
	"path/filepath"
	"strings"
)

type DiagInfo struct {
	DiagName              string
	WarpConectionStatus   bool
	InstalledVersion      string
	PlatformType          string
	TeamName              string
	TeamDomain            string
	SplitTunnelMode       string
	SplitTunnelList       []string
	DeviceProfile         string
	AssignedIPaddress     string
	WarpMode              string
	FallbackDomains       []string
	AlwaysOn              bool
	SwitchLocked          bool
	WiFiDisabled          bool
	EthernetDisabled      bool
	ResolveVia            string
	OnboardingDialogShown bool
	TeamsAuth             bool
	AutoFallback          bool
	CaptivePortalTimeout  int
	SupportURL            string
	Organization          string
	AllowModeSwitch       bool
	AllowUpdates          bool
	AllowLeaveOrg         bool
}

var Info = DiagInfo{}

var ZipPath string

type FileContent struct {
	Data []byte
}

type FileContentMap map[string]FileContent

func ExtractToMemory(zipPath string) (FileContentMap, error) {

	contents := make(FileContentMap)

	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {

		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		var buf bytes.Buffer
		_, err = io.Copy(&buf, rc)
		if err != nil {
			return nil, err
		}

		contents[file.Name] = FileContent{buf.Bytes()}

	}

	return contents, nil

}

func GetInfo(zipPath string, files FileContentMap) DiagInfo {

	Info.DiagName = filepath.Base(zipPath)

	if content, ok := files["platform.txt"]; ok {
		Info.PlatformType = strings.ToLower(string(content.Data))
	}

	if content, ok := files["warp-settings.txt"]; ok {

		settingsLines := strings.Split(string(content.Data), "\n")

		var splitTunnelStart, fallbackDomainsStart, postFallbackSettings int

		for i, line := range settingsLines {
			if strings.Contains(line, "Exclude mode") || strings.Contains(line, "Include mode") {
				splitTunnelStart = i
				Info.SplitTunnelMode = line
			}
			if strings.Contains(line, "Fallback domains") {
				fallbackDomainsStart = i
			}

			if !strings.HasPrefix(line, "  ") {
				postFallbackSettings = i
			}
			// if statements above determine the sections of the settings file.
			// below actually sets the values.

			if strings.Contains(line, "Always On:") {
				if strings.Contains(line, "true") {
					Info.AlwaysOn = true
					continue
				}
				Info.AlwaysOn = false
				continue
			}
			if strings.Contains(line, "Switch Locked:") {
				if strings.Contains(line, "true") {
					Info.SwitchLocked = true
					continue
				}
				Info.SwitchLocked = false
				continue
			}
			if strings.Contains(line, "Mode:") {
				Info.WarpMode = line
				continue
			}

			if strings.Contains(line, "Disabled for Wifi:") {
				if strings.Contains(line, "true") {
					Info.WiFiDisabled = true
					continue
				}
				Info.WiFiDisabled = false
				continue
			}
			if strings.Contains(line, "Disabled for Ethernet:") {
				if strings.Contains(line, "true") {
					Info.EthernetDisabled = true
					continue
				}
				Info.EthernetDisabled = false
				continue
			}

			if strings.Contains(line, "Resolve via:") {
				Info.ResolveVia = line
				continue
			}

			if strings.Contains(line, "Onboarding:") {
				if strings.Contains(line, "true") {
					Info.OnboardingDialogShown = true
					continue
				}
				Info.OnboardingDialogShown = false
				continue
			}
			if strings.Contains(line, "Daemon Teams Auth:") {
				if strings.Contains(line, "true") {
					Info.TeamsAuth = true
					continue
				}
				Info.TeamsAuth = false
				continue
			}
			if strings.Contains(line, "Disable Auto Fallback:") {
				if strings.Contains(line, "true") {
					Info.AutoFallback = true
					continue
				}
				Info.AutoFallback = false
				continue
			}
			if strings.Contains(line, "Support URL:") {
				Info.SupportURL = line
				continue
			}
			if strings.Contains(line, "Organization:") {
				Info.Organization = line
				continue
			}

			if strings.Contains(line, "Allow Mode Switch:") {
				if strings.Contains(line, "true") {
					Info.AllowModeSwitch = true
					continue
				}
				Info.AllowModeSwitch = false
				continue
			}
			if strings.Contains(line, "Allow Updates:") {
				if strings.Contains(line, "true") {
					Info.AllowUpdates = true
					continue
				}
				Info.AllowUpdates = false
				continue

			}
			if strings.Contains(line, "Allowed to Leave Org:") {
				if strings.Contains(line, "true") {
					Info.AllowLeaveOrg = true
					continue
				}
				Info.AllowLeaveOrg = false
				continue
			}

		}

		for _, line := range settingsLines[splitTunnelStart+1 : fallbackDomainsStart] {
			if strings.HasPrefix(line, "  ") {
				splitTunnelEntry := strings.TrimSpace(line)
				Info.SplitTunnelList = append(Info.SplitTunnelList, splitTunnelEntry)

			}
		}
		for _, line := range settingsLines[fallbackDomainsStart+1 : postFallbackSettings] {
			if strings.HasPrefix(line, "  ") {
				fallbackEntry := strings.TrimSpace(line)
				Info.FallbackDomains = append(Info.FallbackDomains, fallbackEntry)
			}
		}

	}

	if content, ok := files["version.txt"]; ok {

		versionContent := strings.Split(string(content.Data), "\n")
		for _, line := range versionContent {
			if strings.Contains(line, "Version:") {
				Info.InstalledVersion = strings.Split(line, " ")[1]
			}
		}
	}

	return Info
}
