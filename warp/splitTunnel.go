package warp

import (
	"fmt"
	"net"
	"strings"
)

var DefaultExcludedCIDRs = []string{
	"10.0.0.0/8",
	"100.64.0.0/10",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.0.0.0/24",
	"192.168.0.0/16",
	"224.0.0.0/24",
	"240.0.0.0/4",
	"255.255.255.255/32",
	"fe80::/10",
	"fd00::/8",
	"ff01::/16",
	"ff02::/16",
	"ff03::/16",
	"ff04::/16",
	"ff05::/16",
}
var Cidrs []string

func (info ParsedDiag) SplitTunnelCheck() (CheckResult, error) {

	SplitTunnelResult := CheckResult{

		CheckName: "IP Address Split Tunnel Check",
		IssueType: "SPLITTUNNEL",
	}

	// Extract CIDR entries

	for _, line := range info.Settings.SplitTunnelList {

		cidr := strings.Split(line, " ")[0] // Only use the first part of the split line as the CIDR ignores comments
		Cidrs = append(Cidrs, cidr)

	}

	// Check if the IP address is in the CIDR entries
	ip := net.ParseIP(info.Network.WarpNetIPv4)
	//fmt.Println("IP Address:", ip) // Add print statement to check IP address
	isInCIDR := false
	var matchedCIDR string
	for _, cidr := range Cidrs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		//fmt.Println("IP Net:", ipNet) // Add print statement to check IP net
		if ipNet.Contains(ip) {
			isInCIDR = true
			matchedCIDR = cidr
			//fmt.Println("IP matched in CIDR:", matchedCIDR) // Add print statement to check if IP is matched in CIDRs
			break
		}
	}

	mode := info.Settings.SplitTunnelMode
	if (strings.Contains(mode, "Exclude mode") && isInCIDR) || (strings.Contains(mode, "Include mode") && !isInCIDR) {
		SplitTunnelResult.CheckPass = true
	} else {
		SplitTunnelResult.CheckPass = false
	}

	if !SplitTunnelResult.CheckPass {
		SplitTunnelResult.Evidence = fmt.Sprintf("Mode: %s\nAssigned IP: %s, Not Matched in Split tunnel CIDRS ", mode, info.Network.WarpNetIPv4)
	} else {
		SplitTunnelResult.Evidence = fmt.Sprintf("Mode: %s\nAssigned IP: %s, Matched in Split tunnel CIDRS: %s", mode, info.Network.WarpNetIPv4, matchedCIDR)
	}

	return SplitTunnelResult, nil
}

func (info ParsedDiag) DefaultExcludeCheck() (CheckResult, error) {

	DefaultExcludeResult := CheckResult{

		CheckName: "Default Exclude Check",
		IssueType: "EXCLUDE_EDITED",
	}
	// Verify default excluded CIDRs

	if strings.Contains(info.Settings.SplitTunnelMode, "Exclude mode") {
		missingCIDRs, allDefaultCIDRsPresent := VerifyDefaultExcludedCIDRs(Cidrs)
		if !allDefaultCIDRsPresent {

			DefaultExcludeResult.IssueType = "EXCLUDE_EDITED"

			DefaultExcludeResult.CheckPass = false
			missingCIDRStr := strings.Join(missingCIDRs, ", ")
			DefaultExcludeResult.Evidence += fmt.Sprintf("Missing default excluded CIDRs: %s", missingCIDRStr)
		} else {
			DefaultExcludeResult.CheckPass = true
			DefaultExcludeResult.Evidence += "All default excluded CIDRs are present"
		}
	}

	return DefaultExcludeResult, nil
}

func VerifyDefaultExcludedCIDRs(cidrs []string) ([]string, bool) {
	missingCIDRs := make([]string, 0)

	for _, defaultCIDR := range DefaultExcludedCIDRs {
		found := false
		for _, cidr := range cidrs {
			if cidr == defaultCIDR {
				found = true
				break
			}
		}
		if !found {
			missingCIDRs = append(missingCIDRs, defaultCIDR)
		}
	}

	return missingCIDRs, len(missingCIDRs) == 0
}
