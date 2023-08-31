package diag

import (
	"archive/zip"
	"bytes"
	"io"
	"path/filepath"
	"strings"
)

var SaveReport, Offline, Verbose, Debug, Plain bool

var ZipPath string

type CheckResult struct {
	CheckID     string
	CheckName   string
	CheckStatus bool
	IssueType   string
	Evidence    string
}

type LogSearchResult struct {
	Filename     string
	SearchTerm   string
	SearchStatus bool
	IssueType    string
	Evidence     string
}

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

type FileContent struct {
	Data []byte
}

type FileContentMap map[string]FileContent

func ExtractZipToMemory(zipPath string) (FileContentMap, error) {

	contents := make(FileContentMap)

	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {

		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		var buf bytes.Buffer
		_, err = io.Copy(&buf, rc)
		if err != nil {
			return nil, err
		}

		contents[file.Name] = FileContent{buf.Bytes()}

	}

	return contents, nil

}

var Info = DiagInfo{}

func GetInfo(zipPath string, files map[string]FileContent) DiagInfo {

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

	return Info
}
