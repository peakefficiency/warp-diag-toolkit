package data

import (
	"archive/zip"
	"bytes"
	"io"
	"path/filepath"
	"strings"
)

type CheckResult struct {
	CheckID      string
	CheckName    string
	CheckPass    bool
	IssueType    string
	Evidence     string
	ReplyMessage string
}
type Diag struct {
	DiagName         string
	InstalledVersion string
	PlatformType     string
	Settings         ParsedSettings
	Account          ParsedAccount
	Network          ParsedNetwork
}

type ParsedAccount struct {
	AccountType  string
	DeviceID     string
	PublicKey    string
	AccountID    string
	Organization string
}

type ParsedDaemonLog struct {
	DeviceProfile string
}

type ParsedNetwork struct {
	WarpNetIPv4 string
	WarpNetIPv6 string
}

type ParsedSettings struct {
	WarpConectionStatus   bool
	SplitTunnelMode       string
	SplitTunnelList       []string
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
	AllowModeSwitch       bool
	AllowUpdates          bool
	AllowLeaveOrg         bool
}

var Info = Diag{}

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

func GetInfo(zipPath string, files FileContentMap) Diag {

	Info.DiagName = filepath.Base(zipPath)

	if content, ok := files["platform.txt"]; ok {
		Info.PlatformType = strings.ToLower(string(content.Data))
	}

	if content, ok := files["warp-account.txt"]; ok {
		accountLines := strings.Split(string(content.Data), "\n")

		for _, line := range accountLines {

			if strings.Contains(line, "Account type:") {
				Info.Account.AccountType = line
				continue
			}
			if strings.Contains(line, "Device ID:") {
				Info.Account.DeviceID = line
				continue
			}
			if strings.Contains(line, "Public key:") {
				Info.Account.PublicKey = line
				continue
			}
			if strings.Contains(line, "Account ID:") {
				Info.Account.AccountID = line
				continue
			}
			if strings.Contains(line, "Organization:") {
				Info.Account.Organization = line
				continue
			}
		}
	}

	if content, ok := files["warp-settings.txt"]; ok {

		settingsLines := strings.Split(string(content.Data), "\n")

		var splitTunnelStart, fallbackDomainsStart, postFallbackSettings int

		for i, line := range settingsLines {
			if strings.Contains(line, "Exclude mode") || strings.Contains(line, "Include mode") {
				splitTunnelStart = i
				Info.Settings.SplitTunnelMode = line
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
					Info.Settings.AlwaysOn = true
					continue
				}
				Info.Settings.AlwaysOn = false
				continue
			}
			if strings.Contains(line, "Switch Locked:") {
				if strings.Contains(line, "true") {
					Info.Settings.SwitchLocked = true
					continue
				}
				Info.Settings.SwitchLocked = false
				continue
			}
			if strings.Contains(line, "Mode:") {
				Info.Settings.WarpMode = line
				continue
			}

			if strings.Contains(line, "Disabled for Wifi:") {
				if strings.Contains(line, "true") {
					Info.Settings.WiFiDisabled = true
					continue
				}
				Info.Settings.WiFiDisabled = false
				continue
			}
			if strings.Contains(line, "Disabled for Ethernet:") {
				if strings.Contains(line, "true") {
					Info.Settings.EthernetDisabled = true
					continue
				}
				Info.Settings.EthernetDisabled = false
				continue
			}

			if strings.Contains(line, "Resolve via:") {
				Info.Settings.ResolveVia = line
				continue
			}

			if strings.Contains(line, "Onboarding:") {
				if strings.Contains(line, "true") {
					Info.Settings.OnboardingDialogShown = true
					continue
				}
				Info.Settings.OnboardingDialogShown = false
				continue
			}
			if strings.Contains(line, "Daemon Teams Auth:") {
				if strings.Contains(line, "true") {
					Info.Settings.TeamsAuth = true
					continue
				}
				Info.Settings.TeamsAuth = false
				continue
			}
			if strings.Contains(line, "Disable Auto Fallback:") {
				if strings.Contains(line, "true") {
					Info.Settings.AutoFallback = true
					continue
				}
				Info.Settings.AutoFallback = false
				continue
			}
			if strings.Contains(line, "Support URL:") {
				Info.Settings.SupportURL = line
				continue
			}

			if strings.Contains(line, "Allow Mode Switch:") {
				if strings.Contains(line, "true") {
					Info.Settings.AllowModeSwitch = true
					continue
				}
				Info.Settings.AllowModeSwitch = false
				continue
			}
			if strings.Contains(line, "Allow Updates:") {
				if strings.Contains(line, "true") {
					Info.Settings.AllowUpdates = true
					continue
				}
				Info.Settings.AllowUpdates = false
				continue

			}
			if strings.Contains(line, "Allowed to Leave Org:") {
				if strings.Contains(line, "true") {
					Info.Settings.AllowLeaveOrg = true
					continue
				}
				Info.Settings.AllowLeaveOrg = false
				continue
			}

		}

		for _, line := range settingsLines[splitTunnelStart+1 : fallbackDomainsStart] {
			if strings.HasPrefix(line, "  ") {
				splitTunnelEntry := strings.TrimSpace(line)
				Info.Settings.SplitTunnelList = append(Info.Settings.SplitTunnelList, splitTunnelEntry)

			}
		}
		for _, line := range settingsLines[fallbackDomainsStart+1 : postFallbackSettings] {
			if strings.HasPrefix(line, "  ") {
				fallbackEntry := strings.TrimSpace(line)
				Info.Settings.FallbackDomains = append(Info.Settings.FallbackDomains, fallbackEntry)
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
