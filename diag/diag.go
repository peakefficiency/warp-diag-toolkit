package diag

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
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

type Info struct {
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

type ZipContent struct {
	Data []byte
}

func ExtractZipToMemory(zipPath string) (map[string]ZipContent, error) {

	contents := make(map[string]ZipContent)

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

		contents[file.Name] = ZipContent{buf.Bytes()}

	}
	if Debug {
		fmt.Println("Files in zip:")
		fmt.Println()
		for filename := range contents {
			fmt.Println(filename)
		}
		if content, ok := contents["connectivity.txt"]; ok {
			fmt.Println()
			fmt.Println("Debug testing connectivity.txt:")
			fmt.Println(string(content.Data))
		}
	}

	return contents, nil

}
