package info

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/peakefficiency/warp-diag-toolkit/diag"
)

var Info = diag.Info{}

func GetInfo(zipPath string, files map[string]diag.ZipContent) diag.Info {

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
					Info.SplitTunnelMode = strings.Split(line, " ")[0]
				} else if strings.HasPrefix(line, "  ") {
					Info.SplitTunnelList += strings.Split(line, " ")[2] + "\n"

				}

			}
		}

	}

	if diag.Debug {
		fmt.Println("Debug check info read: ")
		fmt.Printf("debug Platform type: %s\n", Info.PlatformType)
		fmt.Printf("debug Split tunnel mode: %s\n", Info.SplitTunnelMode)
		fmt.Printf("debug Split tunnel list: \n%s", Info.SplitTunnelList)
		fmt.Printf("debug Fallback domains: \n%s", Info.FallbackDomains)

	}
	return Info
}
