package warp_test

import (
	"testing"

	"github.com/peakefficiency/warp-diag-toolkit/warp"
	"github.com/stretchr/testify/assert"
)

func TestPrintCheckResult(t *testing.T) {
	t.Parallel()

	warp.GetOrLoadConfig(warp.WdcConfig)
	warp.Plain = true
	result := warp.CheckResult{
		CheckID:   "0",
		CheckName: "Warp Version Check",
		IssueType: "OUTDATED_VERSION",
		Evidence:  "Unable to check Linux version automatically, Please verify via package repo https://pkg.cloudflareclient.com/",
	}
	got, _ := result.PrintCheckResult()

	want := "## Warp Version Check\n### OUTDATED_VERSION\n\"It appears that you are not running the latest version of the chosen release\ntrain. \\nPlease attempt to replicate the error using the latest available version\naccording to the details below. \"\n#\n- Evidence: \nUnable to check Linux version automatically, Please verify via package repo https://pkg.cloudflareclient.com/\n"

	assert.Equal(t, want, got, "print check result error")
}
