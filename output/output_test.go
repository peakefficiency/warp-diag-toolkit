package output_test

import (
	"testing"

	"github.com/peakefficiency/warp-diag-toolkit/data"
	"github.com/peakefficiency/warp-diag-toolkit/output"
	"github.com/stretchr/testify/assert"
)

func TestPrintCheckResult(t *testing.T) {
	result := data.CheckResult{
		CheckID:   "0",
		CheckName: "Warp Version Check",
		IssueType: "OUTDATED_VERSION",
		Evidence:  "Unable to check Linux version automatically, Please verify via package repo https://pkg.cloudflareclient.com/",
	}
	got, _ := output.PrintCheckResult(result)

	want := "message "

	assert.Equal(t, want, got, "print check result error")
}
