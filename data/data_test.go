package data_test

import (
	"archive/zip"
	"bytes"
	"os"
	"testing"

	"github.com/peakefficiency/warp-diag-toolkit/data"
	"github.com/stretchr/testify/assert"
)

func TestZipToInfo(t *testing.T) {
	//will fail if tests not parellel as only one diag to be processed at a time
	t.Parallel()

	realZipPath := "testdata/warp-debugging-info-20230831-185328.zip"
	files, err := data.ExtractToMemory(realZipPath)
	if err != nil {
		t.Error("Some error extracting zip", err)
	}

	info := data.GetInfo(realZipPath, files)

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
	files := data.FileContentMap{
		"platform.txt": data.FileContent{
			Data: []byte("windows"),
		},
		"warp-settings.txt": data.FileContent{
			Data: []byte("Exclude mode\n  192.168.1.10\n  192.168.1.20\n"),
		},
	}

	info := data.GetInfo(zipPath, files)

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
	emptyfiles := data.FileContentMap{}

	invalidinfo := data.GetInfo(zipPath, emptyfiles)

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

func createTestZipFile() (string, error) {
	zipFilePath := "test_data.zip"
	file, err := os.Create(zipFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)

	// Add test files to the zip
	file1, err := zipWriter.Create("file1.txt")
	if err != nil {
		return "", err
	}
	file1.Write([]byte("This is file 1"))

	file2, err := zipWriter.Create("file2.txt")
	if err != nil {
		return "", err
	}
	file2.Write([]byte("This is file 2"))

	zipWriter.Close()

	return zipFilePath, nil
}

func TestExtractZipToMemory(t *testing.T) {
	t.Parallel()
	// Create a test zip file
	zipFilePath, err := createTestZipFile()
	if err != nil {
		t.Errorf("Error creating test zip file: %v", err)
		return
	}
	defer os.Remove(zipFilePath)

	contents, err := data.ExtractToMemory(zipFilePath)
	if err != nil {
		t.Errorf("Error extracting zip file: %v", err)
		return
	}

	// Check the extracted contents
	expectedFile1Content := []byte("This is file 1")
	expectedFile2Content := []byte("This is file 2")

	file1Content, ok := contents["file1.txt"]
	if !ok {
		t.Error("Expected file1.txt to be extracted")
		return
	}

	file2Content, ok := contents["file2.txt"]
	if !ok {
		t.Error("Expected file2.txt to be extracted")
		return
	}

	if !bytes.Equal(expectedFile1Content, file1Content.Data) {
		t.Error("file1.txt content does not match expected")
		return
	}

	if !bytes.Equal(expectedFile2Content, file2Content.Data) {
		t.Error("file2.txt content does not match expected")
		return
	}

	//check invalid file and content
	_, ok = contents["file3.txt"]
	if ok {
		t.Fatal("expected file not found but received 'ok'")
		return
	}

}

func TestExtractToMemoryRealFile(t *testing.T) {
	t.Parallel()

	contents, err := data.ExtractToMemory("testdata/test.zip")
	if err != nil {
		t.Errorf("Error extracting zip file: %v", err)
		return
	}

	// Check the extracted contents
	expectedFile1Content := []byte("This is real file 1")
	expectedFile2Content := []byte("This is real file 2")

	file1Content, ok := contents["file1.txt"]
	if !ok {
		t.Error("Expected file1.txt to be extracted")
		return
	}

	file2Content, ok := contents["file2.txt"]
	if !ok {
		t.Error("Expected file2.txt to be extracted")
		return
	}

	if !bytes.Equal(expectedFile1Content, file1Content.Data) {
		t.Error("file1.txt content does not match expected")
		return
	}

	if !bytes.Equal(expectedFile2Content, file2Content.Data) {
		t.Error("file2.txt content does not match expected")
		return
	}

	//check invalid file and content
	_, ok = contents["file3.txt"]
	if ok {
		t.Fatal("expected file not found but received 'ok'")
		return
	}

}
