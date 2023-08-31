package info

import (
	"path/filepath"
	"strings"

	"github.com/peakefficiency/warp-diag-toolkit/diag"
)

type DiagInfo struct {
	DiagName            string
	WarpConectionStatus bool
	InstalledVersion    string
	PlatformType        string
	TeamName            string
	TeamDomain          string
	SplitTunnelMode     string
	SplitTunnelList     string
	DeviceProfile       string
	AssignedIPaddress   string
	WarpMode            string
	FallbackDomains     string
}

var Info = DiagInfo{}

func GetInfo(zipPath string, files diag.FileContentMap) DiagInfo {

	Info.DiagName = filepath.Base(zipPath)

	for name, content := range files {
		if name == "platform.txt" {
			Info.PlatformType = strings.ToLower(string(content.Data))
		}

		if name == "warp-settings.txt" {

			lines := strings.Split(string(content.Data), "\n")

			// Parse split tunnel list
			for _, line := range lines {
				if strings.Contains(line, "Exclude mode") || strings.Contains(line, "Include mode") && !strings.Contains(line, "Fallback domains") {
					Info.SplitTunnelMode = line
				} else if strings.HasPrefix(line, "  ") {
					Info.SplitTunnelList += strings.Split(line, " ")[2] + "\n"

				}

			}
		}

	}

	return Info
}
