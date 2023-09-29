package warp

import (
	"fmt"
	"strings"
)

type LogSearchResult struct {
	Filename     string
	SearchTerm   string
	SearchStatus bool
	IssueType    string
	Evidence     string
}

var LogSearchOutput = map[string]LogSearchResult{}

func LogSearch(contents map[string]FileContent) map[string]LogSearchResult {
	// search logic

	for _, logPattern := range Conf.LogPatternsByIssue {

		searchFilename := logPattern.SearchFile
		if Info.PlatformType == "windows" && searchFilename == "ps.txt" {
			searchFilename = "processes.txt"
		}

		content, found := contents[searchFilename]
		if !found {
			continue
		}

		fileContent := string(content.Data)

		for issueType, issue := range logPattern.Issue {

			evidence := []string{}

			for _, searchTerm := range issue.SearchTerms {

				for _, line := range strings.Split(fileContent, "\n") {

					if strings.Contains(line, searchTerm) {
						evidence = append(evidence, line)
					}

				}
			}

			if len(evidence) > 0 {
				LogSearchOutput[issueType] = LogSearchResult{
					IssueType: issueType,
					Evidence:  strings.Join(evidence, "\n"),
				}
			}

		}

	}
	if Debug {
		fmt.Println("Log Search Output:")
		for issueType, result := range LogSearchOutput {
			fmt.Printf(" IssueType: %s\n", issueType)
			fmt.Printf(" Evidence:\n%s\n", result.Evidence)
		}
	}
	return LogSearchOutput
}
