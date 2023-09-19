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

	assert.Equal(t, true, info.AlwaysOn, "always on not detected correctly")
	//needs some work to define the test elegantly then can fix implementation

	expectedSplitTunnelIPs := []string{
		"10.0.0.0/8",
		"100.64.0.0/10",
		"169.254.0.0/16 (DHCP Unspecified)",
		"172.16.0.0/12",
		"192.0.0.0/24",
		"192.168.0.0/16",
		"224.0.0.0/24",
		"240.0.0.0/4",
		"255.255.255.255/32 (DHCP Broadcast)",
		"fe80::/10 (IPv6 Link Local)",
		"fd00::/8",
		"ff01::/16",
		"ff02::/16",
		"ff03::/16",
		"ff04::/16",
		"ff05::/16",
		"*.wikipedia.org",
		"*.en.wikipedia.org",
	}

	assert.Equal(t, expectedSplitTunnelIPs, info.SplitTunnelList, "Split tunnel list doesnt match")
	//expectedFallbackDomains := []string{
	//	"home.arpa",
	//	"wikipedia.org	-> [8.8.8.8]",
	//	"intranet",
	//	"internal",
	//	"private",
	//	"localdomain",
	//	"domain",
	//	"lan",
	//	"home",
	//	"host",
	//	"corp",
	//	"local",
	//	"localhost",
	//	"invalid",
	//	"test",
	//}
	//assert.Equal(t, expectedFallbackDomains, info.FallbackDomains, "Fallback domains dont match")

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
	if len(invalidinfo.SplitTunnelList) != 0 {
		t.Errorf("Expected SplitTunnelList to be empty, got %v", invalidinfo.SplitTunnelList)
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
