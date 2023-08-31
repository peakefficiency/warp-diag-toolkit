package info_test

import (
	"testing"

	"github.com/peakefficiency/warp-diag-toolkit/diag"
	"github.com/peakefficiency/warp-diag-toolkit/info"
	"github.com/stretchr/testify/assert"
)

func TestZipToInfo(t *testing.T) {
	//will fail if tests not parellel as only one diag to be processed at a time
	t.Parallel()

	realZipPath := "testdata/warp-debugging-info-20230831-185328.zip"
	files, err := diag.ExtractToMemory(realZipPath)
	if err != nil {
		t.Error("Some error extracting zip", err)
	}

	info := info.GetInfo(realZipPath, files)

	if info.DiagName != "warp-debugging-info-20230831-185328.zip" {
		t.Errorf("Expected DiagName to be %s, got %s", "warp-debugging-info-20230831-185328.zip", info.DiagName)
	}
	if info.PlatformType != "macos" {
		t.Errorf("Expected PlatformType to be %s, got %s", "macos", info.PlatformType)
	}
	assert.Containsf(t, info.SplitTunnelMode, "Exclude", "expected Split Tunne mode to be Exclude got %s", info.SplitTunnelMode)

	//needs some work to define the test elegantly then can fix implementation
	//assert.Equalf(t, info.SplitTunnelList, "192.168.1.10\n192.168.1.20\n10.0.0.0/8\n100.64.0.0/10\n169.254.0.0/16\n172.16.0.0/12\n192.0.0.0/24\n192.168.0.0/16\n224.0.0.0/24\n240.0.0.0/4\n255.255.255.255/32\nfe80::/10\nfd00::/8\nff01::/16\nff02::/16\nff03::/16\nff04::/16\nff05::/16\n*.wikipedia.org\n*.en.wikipedia.org\nhome.arpa\nwikipedia.org\nintranet\ninternal\nprivate\nlocaldomain\ndomain\nlan\nhome\nhost\ncorp\nlocal\nlocalhost\ninvalid\ntest", "Expected SplitTunnelList to be %s, got %s", "192.168.1.10\n192.168.1.20\n10.0.0.0/8\n100.64.0.0/10\n169.254.0.0/16\n172.16.0.0/12\n192.0.0.0/24\n192.168.0.0/16\n224.0.0.0/24\n240.0.0.0/4\n255.255.255.255/32\nfe80::/10\nfd00::/8\nff01::/16\nff02::/16\nff03::/16\nff04::/16\nff05::/16\n*.wikipedia.org\n*.en.wikipedia.org\nhome.arpa\nwikipedia.org\nintranet\ninternal\nprivate\nlocaldomain\ndomain\nlan\nhome\nhost\ncorp\nlocal\nlocalhost\ninvalid\ntest", info.SplitTunnelList)
}

func TestGetInfo(t *testing.T) {
	//will fail if tests not parellel as only one diag to be processed at a time
	t.Parallel()

	zipPath := "/path/to/zipfile"
	files := diag.FileContentMap{
		"platform.txt": diag.FileContent{
			Data: []byte("windows"),
		},
		"warp-settings.txt": diag.FileContent{
			Data: []byte("Exclude mode\n  192.168.1.10\n  192.168.1.20\n"),
		},
	}

	info := info.GetInfo(zipPath, files)

	if info.DiagName != "zipfile" {
		t.Errorf("Expected DiagName to be %s, got %s", "zipfile", info.DiagName)
	}
	if info.PlatformType != "windows" {
		t.Errorf("Expected PlatformType to be %s, got %s", "windows", info.PlatformType)
	}
	if info.SplitTunnelMode != "Exclude mode" {
		t.Errorf("Expected SplitTunnelMode to be %s, got %s", "Exclude", info.SplitTunnelMode)
	}
	if info.SplitTunnelList != "192.168.1.10\n192.168.1.20\n" {
		t.Errorf("Expected SplitTunnelList to be %s, got %s", "192.168.1.10\n192.168.1.20\n", info.SplitTunnelList)
	}
}

func TestGetInfoEmptyFiles(t *testing.T) {
	t.Parallel()

	zipPath := "/path/to/zipfile"
	emptyfiles := diag.FileContentMap{}

	invalidinfo := info.GetInfo(zipPath, emptyfiles)

	if invalidinfo.DiagName != "zipfile" {
		t.Errorf("Expected DiagName to be %s, got %s", "zipfile", invalidinfo.DiagName)
	}
	if invalidinfo.PlatformType != "" {
		t.Errorf("Expected PlatformType to be %s, got %s", "", invalidinfo.PlatformType)
	}
	if invalidinfo.SplitTunnelMode != "" {
		t.Errorf("Expected SplitTunnelMode to be %s, got %s", "", invalidinfo.SplitTunnelMode)
	}
	if invalidinfo.SplitTunnelList != "" {
		t.Errorf("Expected SplitTunnelList to be %s, got %s", "", invalidinfo.SplitTunnelList)
	}
}
