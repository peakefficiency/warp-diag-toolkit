package checks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/peakefficiency/warp-diag-toolkit/data"
)

const (
	macReleaseURL          = "https://install.appcenter.ms/api/v0.1/apps/cloudflare/1.1.1.1-macos-1/distribution_groups/release/public_releases?scope=tester"
	macBetaURL             = "https://install.appcenter.ms/api/v0.1/apps/cloudflare/1.1.1.1-macos/distribution_groups/beta/public_releases?scope=tester"
	windowsReleaseURL      = "https://install.appcenter.ms/api/v0.1/apps/cloudflare/1.1.1.1-windows-1/distribution_groups/release/public_releases?scope=tester"
	windowsBetaURL         = "https://install.appcenter.ms/api/v0.1/apps/cloudflare/1.1.1.1-windows/distribution_groups/beta/public_releases?scope=tester"
	linuxReleaseURL        = "https://pkg.cloudflareclient.com/"
	windowsDownloadURL     = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-windows-1/distribution_groups/release"
	windowsBetaDownloadURL = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-windows/distribution_groups/beta"
	macDownloadURL         = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-macos-1/distribution_groups/release"
	macBetaDownloadURL     = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-macos/distribution_groups/beta"
)

type Release struct {
	ID           int    `json:"id"`
	ShortVersion string `json:"short_version"`
	Version      string `json:"version"`
}

func fetchLatestVersionURL(beta bool) (string, error) {
	switch data.Info.PlatformType {
	case "windows":
		if beta {
			return windowsBetaURL, nil
		}
		return windowsReleaseURL, nil
	case "mac":
		if beta {
			return macBetaURL, nil
		}
		return macReleaseURL, nil
	case "linux":
		return linuxReleaseURL, nil
	default:
		return "", fmt.Errorf("unknown platform type")
	}
}

func fetchLatestVersion(beta bool) (string, error) {

	if data.Info.PlatformType == "linux" {
		return "", nil // Do not fetch the latest version for Linux
	}

	url, err := fetchLatestVersionURL(beta)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1 Safari/605.1.15")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest version: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var releases []Release
	err = json.Unmarshal(bodyBytes, &releases)
	if err != nil {
		return "", fmt.Errorf("failed to decode JSON response: %v", err)
	}

	if len(releases) == 0 {
		return "", fmt.Errorf("no releases found")
	}

	if data.Info.PlatformType == "windows" {
		return releases[0].Version, nil
	} else {
		return releases[0].ShortVersion, nil
	}
}

func RunVersionCheck() data.CheckResult {

	INSTALLEDversion, err := version.NewVersion(data.Info.InstalledVersion)
	if err != nil {
		fmt.Println("Error fetching installed version:", err)
	}

	var latestVersion string
	var latestBetaVersion string

	latestVersion, err = fetchLatestVersion(false)
	if err != nil {
		fmt.Println("Error fetching latest version:", err)
	}
	RELEASEversion, err := version.NewVersion(latestVersion)
	if err != nil {
		fmt.Println("Error fetching release version:", err)
	}
	if err != nil && platform != types.LinuxPlatform {
		return types.CheckResult{
			CheckName: "Warp version verification",
			Details:   fmt.Sprintf("Failed to fetch latest version: %v", err),
		}
	}

	isBeta := false

	if platform == types.LinuxPlatform {
		linuxCheckResult := types.CheckResult{
			CheckID:   "0",
			CheckName: "Warp version verification",
			Success:   false,
			Details:   fmt.Sprintf("Installed version: `%s`\nPlease ensure you are using the latest version.\nVerify the latest version here(%s)", installedVersion, linuxReleaseURL),
		}
		if debug {
			fmt.Printf("Linux Check Result: %+v\n", linuxCheckResult)
		}
		return linuxCheckResult
	}

	if INSTALLEDversion.GreaterThan(RELEASEversion) {
		isBeta = true
		latestBetaVersion, err = fetchLatestVersion(platform, true)
		if err != nil && platform != types.LinuxPlatform {
			return types.CheckResult{
				CheckName: "Warp version verification",
				Details:   fmt.Sprintf("Failed to fetch latest beta version: %v", err),
			}
		}
	}

	var downloadURL string
	if isBeta {
		switch platform {
		case types.WindowsPlatform:
			downloadURL = windowsBetaDownloadURL
		case types.MacPlatform:
			downloadURL = macBetaDownloadURL
		}
	} else {
		switch platform {
		case types.WindowsPlatform:
			downloadURL = windowsDownloadURL
		case types.MacPlatform:
			downloadURL = macDownloadURL
		}
	}

	for _, badVersion := range config.BadVersions {
		if installedVersion == badVersion {
			return types.CheckResult{
				CheckID:   "0",
				CheckName: "Warp version verification",
				Success:   false,
				Details:   "Issue: BADVERSION\n" + fmt.Sprintf("Installed version `%s` is a known bad version", installedVersion),
			}
		}
	}

	if isBeta {
		if INSTALLEDversion.Equal(RELEASEversion) {
			return types.CheckResult{
				CheckID:   "0",
				CheckName: "Warp version verification",
				Success:   true,
				Details:   fmt.Sprintf("Installed beta version: `%s`.\nLatest stable version: `%s`", installedVersion, latestVersion),
			}
		} else {
			return types.CheckResult{
				CheckID:   "0",
				CheckName: "Warp version verification",
				Success:   false,
				Details:   fmt.Sprintf("Issue: OUTDATED_VERSION\nInstalled version `%s` appears to be a beta, but not the latest beta %s.\nLatest stable version: %s.\n\nDownload the latest beta version here(%s)", installedVersion, latestBetaVersion, latestVersion, downloadURL),
			}
		}
	} else {
		if INSTALLEDversion.Equal(RELEASEversion) {
			return types.CheckResult{
				CheckID:   "0",
				CheckName: "Warp version verification",
				Success:   true,
				Details:   fmt.Sprintf("Installed version: %s, Latest version: %s", installedVersion, latestVersion),
			}
		} else {
			return types.CheckResult{
				CheckID:   "0",
				CheckName: "Warp version verification",
				Success:   false,
				Details:   fmt.Sprintf("Issue: OUTDATED_VERSION\nVersion mismatch: Installed: %s, Latest: %s.\n\nDownload the latest version here:(%s)", installedVersion, latestVersion, downloadURL),
			}
		}
	}
}
