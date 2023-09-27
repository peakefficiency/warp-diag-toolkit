package checks_test

import (
	"testing"

	"github.com/peakefficiency/warp-diag-toolkit/checks"
	"github.com/peakefficiency/warp-diag-toolkit/data"
	"github.com/stretchr/testify/assert"
)

// running tests with live external endpoints for now - testing not parellel to help prevent rate limiting and incorret data
func TestGetLatestReleaseVersionMac(t *testing.T) {

	releaseWant := "2023.7.159.0"
	betaWant := "2023.9.109.1"
	macVersions, err := checks.LatestMacVersions()
	releaseGot := macVersions.Release
	betaGot := macVersions.Beta

	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, betaWant, betaGot, "beta version fail")
	assert.Equal(t, releaseWant, releaseGot, "Release version fail")
}

func TestGetLatestVersionsWindows(t *testing.T) {

	betaWinWant := "2023.9.107.1"
	releaseWinWant := "2023.7.160.0"
	winVersions, err := checks.LatestWinVersions()

	betaWinGot := winVersions.Beta
	releaseWinGot := winVersions.Release
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, betaWinWant, betaWinGot, "beta version fail")
	assert.Equal(t, releaseWinWant, releaseWinGot, "Release version fail")

}

func TestVersionCheckLinux(t *testing.T) {

	data.Info.PlatformType = "linux"
	got := checks.VersionCheck()

	want := data.CheckResult{
		CheckID:   "0",
		CheckName: "Warp Version Check",
		IssueType: "OUTDATED_VERSION",
		Evidence:  "Unable to check Linux version automatically, Please verify via package repo https://pkg.cloudflareclient.com/",
	}

	assert.Equal(t, want, got, "Output doesnt match")

}

func TestVersionWindowsOldRelease(t *testing.T) {

	data.Info.PlatformType = "windows"
	data.Info.InstalledVersion = "2023.7.100.0"

	got := checks.VersionCheck()

	want := data.CheckResult{
		CheckID:     "0",
		CheckName:   "Warp Version Check",
		IssueType:   "OUTDATED_VERSION",
		Evidence:    "installed version: 2023.7.100.0, Latest Release version: 2023.7.160.0",
		CheckStatus: false,
	}

	assert.Equal(t, want, got, "Check ID fail")

}

func TestVersionWindowsOldBeta(t *testing.T) {

	data.Info.PlatformType = "windows"
	data.Info.InstalledVersion = "2023.7.200.0"

	got := checks.VersionCheck()

	want := data.CheckResult{
		CheckID:     "0",
		CheckName:   "Warp Version Check",
		IssueType:   "OUTDATED_VERSION",
		Evidence:    "installed version: 2023.7.200.0, Which appears to be a beta as it is newer than the latest release: 2023.7.160.0,  but not the latest beta which is: 2023.9.107.1",
		CheckStatus: false,
	}

	assert.Equal(t, want, got, "Check ID fail")
}

func TestVersionMacOldRelease(t *testing.T) {
	data.Info.PlatformType = "mac"
	data.Info.InstalledVersion = "2023.7.100.0"

	got := checks.VersionCheck()

	want := data.CheckResult{
		CheckID:     "0",
		CheckName:   "Warp Version Check",
		IssueType:   "OUTDATED_VERSION",
		Evidence:    "installed version: 2023.7.100.0, Latest Release version: 2023.7.159.0",
		CheckStatus: false,
	}

	assert.Equal(t, want, got, "Check ID fail")

}

func TestVersionMacOldBeta(t *testing.T) {
	data.Info.PlatformType = "mac"
	data.Info.InstalledVersion = "2023.7.200.0"

	got := checks.VersionCheck()

	want := data.CheckResult{
		CheckID:     "0",
		CheckName:   "Warp Version Check",
		IssueType:   "OUTDATED_VERSION",
		Evidence:    "installed version: 2023.7.200.0, Which appears to be a beta as it is newer than the latest release: 2023.7.159.0,  but not the latest beta which is: 2023.9.109.1",
		CheckStatus: false,
	}

	assert.Equal(t, want, got, "Check ID fail")
}
