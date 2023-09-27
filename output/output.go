package output

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/peakefficiency/warp-diag-toolkit/checks"
	"github.com/peakefficiency/warp-diag-toolkit/cli"
	"github.com/peakefficiency/warp-diag-toolkit/config"
	"github.com/peakefficiency/warp-diag-toolkit/data"
)

func DumpFiles(files data.FileContentMap, filename string) {

	if filename != "" {
		// Dump specific file
		if content, ok := files[filename]; ok {
			fmt.Println(filename)
			fmt.Println(string(content.Data))
		} else {
			fmt.Printf("File %s not found in zip\n", filename)
		}
	} else {
		// Dump all files
		fmt.Println("# File Contents")

		for name, content := range files {
			fmt.Printf("## %s\n", name)
			fmt.Println(string(content.Data))
		}
	}

}

func ReportInfo(info data.Diag) (string, error) {
	var markdown strings.Builder

	markdown.WriteString("## Warp Diag Information\n")

	markdown.WriteString(fmt.Sprintf("* Name: %s\n", info.DiagName))
	markdown.WriteString(fmt.Sprintf("* Platform: %s\n", info.PlatformType))

	if cli.Plain {
		return markdown.String(), nil
	}

	return glamour.Render(markdown.String(), "dark")
}

func ReportLogSearch(results map[string]checks.LogSearchResult) (string, error) {
	var markdown strings.Builder

	markdown.WriteString("## Log Search Results\n")

	for issueType, result := range results {
		reply := config.Conf.ReplyByIssueType[issueType]

		markdown.WriteString(fmt.Sprintf("### %s\n", issueType))
		markdown.WriteString(fmt.Sprintf("%s\n", reply.Message))
		markdown.WriteString(fmt.Sprintf("- Evidence: \n%s\n", result.Evidence))
	}

	if cli.Plain {
		return markdown.String(), nil
	}

	return glamour.Render(markdown.String(), "dark")
}

func PrintCheckResult(result data.CheckResult) (string, error) {
	var markdown strings.Builder

	if !result.CheckPass {

	}

	reply := config.Conf.ReplyByIssueType["result.IssueType"]

	markdown.WriteString(fmt.Sprintf("## %s\n", result.CheckName))

	markdown.WriteString(fmt.Sprintf("### %s\n", result.IssueType))
	markdown.WriteString(fmt.Sprintf("%s\n", reply.Message))
	markdown.WriteString(fmt.Sprintf("- Evidence: \n%s\n", result.Evidence))

	if cli.Plain {
		return markdown.String(), nil
	}

	return glamour.Render(markdown.String(), "dark")
}
