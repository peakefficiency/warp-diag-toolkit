package warp_test

import (
	"testing"

	"github.com/peakefficiency/warp-diag-toolkit/warp"
	"github.com/stretchr/testify/assert"
)

// running tests with live external endpoints for now - testing not parellel to help prevent rate limiting and incorret data
func TestGetLatestReleaseVersionMac(t *testing.T) {
	t.Parallel()

	releaseWant := "2023.9.252.0"
	betaWant := "2023.9.109.1"
	macVersions, err := warp.LatestMacVersions()
	releaseGot := macVersions.Release
	betaGot := macVersions.Beta

	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, betaWant, betaGot, "beta version fail")
	assert.Equal(t, releaseWant, releaseGot, "Release version fail")
}

func TestGetLatestVersionsWindows(t *testing.T) {
	t.Parallel()
	betaWinWant := "2023.9.107.1"
	releaseWinWant := "2023.9.248.0"
	winVersions, err := warp.LatestWinVersions()

	betaWinGot := winVersions.Beta
	releaseWinGot := winVersions.Release
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, betaWinWant, betaWinGot, "beta version fail")
	assert.Equal(t, releaseWinWant, releaseWinGot, "Release version fail")

}

func TestVersionCheckLinux(t *testing.T) {
	t.Parallel()

	info := warp.ParsedDiag{}

	info.PlatformType = "linux"

	got := info.VersionCheck()

	want := warp.CheckResult{
		CheckID:   "0",
		CheckName: "Warp Version Check",
		IssueType: "OUTDATED_VERSION",
		Evidence:  "Unable to check Linux version automatically, Please verify via package repo https://pkg.cloudflareclient.com/",
	}

	assert.Equal(t, want, got, "Output doesnt match")

}

func TestVersionWindowsOldRelease(t *testing.T) {
	t.Parallel()
	info := warp.ParsedDiag{}

	info.PlatformType = "windows"
	info.InstalledVersion = "2023.7.100.0"

	got := info.VersionCheck()

	want := warp.CheckResult{
		CheckID:   "0",
		CheckName: "Warp Version Check",
		IssueType: "OUTDATED_VERSION",
		Evidence:  "installed version: 2023.7.100.0, Latest Release version: 2023.9.248.0",
		CheckPass: false,
	}

	assert.Equal(t, want, got, "Check ID fail")

}

//func TestVersionWindowsOldBeta(t *testing.T) {

//	warp.Info.PlatformType = "windows"
//	warp.Info.InstalledVersion = "2023.7.200.0"

//	got := warp.VersionCheck()

//	want := warp.CheckResult{
//		CheckID:   "0",
//		CheckName: "Warp Version Check",
//		IssueType: "OUTDATED_VERSION",
//		Evidence:  "installed version: 2023.7.200.0, Which appears to be a beta as it is newer than the latest release: 2023.9.248.0,  but not the latest beta which is: 2023.9.107.1",
//		CheckPass: false,
//	}

// assert.Equal(t, want, got, "Check ID fail")
//}

func TestVersionMacOldRelease(t *testing.T) {

	info := warp.ParsedDiag{}

	info.PlatformType = "mac"
	info.InstalledVersion = "2023.7.100.0"

	got := info.VersionCheck()

	want := warp.CheckResult{
		CheckID:   "0",
		CheckName: "Warp Version Check",
		IssueType: "OUTDATED_VERSION",
		Evidence:  "installed version: 2023.7.100.0, Latest Release version: 2023.9.252.0",
		CheckPass: false,
	}

	assert.Equal(t, want, got, "Check ID fail")

}

//while old beta not updated to be newer than release
//func TestVersionMacOldBeta(t *testing.T) {
//	warp.Info.PlatformType = "mac"
//	warp.Info.InstalledVersion = "2023.7.200.0"
//
//	got := warp.VersionCheck()
//
//	want := warp.CheckResult{
//		CheckID:   "0",
//		CheckName: "Warp Version Check",
//		IssueType: "OUTDATED_VERSION",
//		Evidence:  "installed version: 2023.7.200.0, Which appears to be a beta as it is newer than the latest release: 2023.9.252.0,  but not the latest beta which is: 2023.9.109.1",
//		CheckPass: false,
//	}
//
//	assert.Equal(t, want, got, "Check ID fail")//
//}
