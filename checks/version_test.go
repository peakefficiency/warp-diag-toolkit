package checks_test

import (
	"testing"

	"github.com/peakefficiency/warp-diag-toolkit/checks"
	"github.com/peakefficiency/warp-diag-toolkit/data"
)

// testing not parrelle
func TestGetLatestReleaseVersionMac(t *testing.T) {

	data.Info.PlatformType = "mac"
	want := "2023.7.159.0"
	got, err := checks.FetchLatestVersionFrom(checks.MacReleaseURL)
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestGetLatestBetaVersionWindows(t *testing.T) {

	data.Info.PlatformType = "windows"
	want := "2023.9.107.1"
	got, err := checks.FetchLatestVersionFrom(checks.WindowsBetaURL)
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
