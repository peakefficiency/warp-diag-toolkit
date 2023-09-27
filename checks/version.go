package checks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/peakefficiency/warp-diag-toolkit/data"
)

var VersionCheckResult = data.CheckResult{}

const (
	MacReleaseURL          = "https://install.appcenter.ms/api/v0.1/apps/cloudflare/1.1.1.1-macos-1/distribution_groups/release/public_releases?scope=tester"
	MacBetaURL             = "https://install.appcenter.ms/api/v0.1/apps/cloudflare/1.1.1.1-macos/distribution_groups/beta/public_releases?scope=tester"
	WindowsReleaseURL      = "https://install.appcenter.ms/api/v0.1/apps/cloudflare/1.1.1.1-windows-1/distribution_groups/release/public_releases?scope=tester"
	WindowsBetaURL         = "https://install.appcenter.ms/api/v0.1/apps/cloudflare/1.1.1.1-windows/distribution_groups/beta/public_releases?scope=tester"
	LinuxPKGurl            = "https://pkg.cloudflareclient.com/"
	WindowsDownloadURL     = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-windows-1/distribution_groups/release"
	WindowsBetaDownloadURL = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-windows/distribution_groups/beta"
	MacDownloadURL         = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-macos-1/distribution_groups/release"
	MacBetaDownloadURL     = "https://install.appcenter.ms/orgs/cloudflare/apps/1.1.1.1-macos/distribution_groups/beta"
)

const (
	ForBeta    = true
	ForRelease = false
)

type Release struct {
	ID              int       `json:"id"`
	ShortVersion    string    `json:"short_version"`
	Version         string    `json:"version"`
	UploadedAt      time.Time `json:"uploaded_at"`
	MandatoryUpdate bool      `json:"mandatory_update"`
	Enabled         bool      `json:"enabled"`
}

//helper function http call to get latest release which is the first json object in the response.

func FetchLatestVersionFrom(url string) (string, error) {

	client := &http.Client{
		Timeout: time.Second * 1,
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
