package kernel

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alexeyco/simpletable"
)

const (
	colorCyan    = "\033[36m"
	colorReset   = "\033[0m"
	colorDefault = "\x1b[39m"
	colorRed     = "\x1b[91m"
)

type checkResult struct {
	Filename   string
	Name       string
	Value      string
	Configured bool
}

func Execute(clusterType, nodeType int) {
	// cluster - 0 - SNC, 1 - 2x2 small, 2 - 2x2 regular, 3 - 4x2
	// node - 0 for Core, 1 for Data, 2 for SNC (default)
	allNodesKernel()

	if clusterType == 0 || nodeType == 0 {
		onlyCoreAndSnc()
	}
	if clusterType == 0 || nodeType == 1 {
		onlyDataAndSnc()
	}
}

func allNodesKernel() {
	l1 := "\n\nThe following kernel settings are required for both Core & Data Nodes, also for SNC:\n"
	l2 := "ATTENTION! You will need to create /etc/systemd/system illumio-pce.service.d/override.conf file after PCE RPM installed.\n* For CentOS/RHEL 7.x/8.x with systemd.\n* Remember reload the daemon:"
	l3 := "systemctl daemon-reload && systemctl restart illumio-pce.service"
	l4 := "\n* Verify it is in effect: "
	l5 := "sudo -u ilo-pce systemctl show illumio-pce.service | egrep 'LimitCORE|LimitNPROC|LimitNOFILE'"
	l6 := "\n* If the environment prefers init.d rather than systemd, plese check the Illumio manual on how to set init.d value for PCE."

	fmt.Println(string(colorCyan), l1, string(colorReset), l2, string(colorCyan), l3, string(colorReset), l4, string(colorCyan), l5, string(colorReset), l6)

	settings := []string{
		"LimitCORE=0", "LimitNOFILE=65535", "LimitNPROC=65535",
	}

	filePath := "/etc/systemd/system/illumio-pce.service.d/override.conf"
	output := grepFiles(settings, filePath)
	createTable(output)
}

func onlyCoreAndSnc() {
	l1 := "\n\nThe following kernel settings are required for Core Nodes and SNC only:\n"
	l2 := "If your settings are greater than these, you do not need to lower them.\n* Remember to apply the change:"
	l3 := "sysctl -p /etc/sysctl.d/99-illumio.conf"

	fmt.Println(string(colorCyan), l1, string(colorReset), l2, string(colorCyan), l3, string(colorReset))

	settings := []string{
		"fs.file-max = 2000000", "net.core.somaxconn = 16384", "net.nf_conntrack_max = 1048576",
	}

	filePath := "/etc/sysctl.d/99-illumio.conf"
	output := grepFiles(settings, filePath)
	createTable(output)

	// kernel module check
	l4 := "\n\nThe following kernel conntrack module settings are required for Core Nodes and SNC only:\n"
	l5 := "If your settings are greater than these, you do not need to lower them.\n* Remember to enable first:"
	l6 := "modprobe nf_conntrack"

	fmt.Println(string(colorCyan), l4, string(colorReset), l5, string(colorCyan), l6, string(colorReset))

	settings = []string{
		"options nf_conntrack hashsize=262144",
	}

	filePath = "/etc/modprobe.d/illumio.conf"
	output = grepFiles(settings, filePath)
	createTable(output)
}

func onlyDataAndSnc() {
	l1 := "\n\nThe following kernel settings are required for Data Nodes and SNC only:\n"
	l2 := "If your settings are greater than these, you do not need to lower them.\n* Remember to apply the change:"
	l3 := "sysctl -p /etc/sysctl.d/99-illumio.conf"

	fmt.Println(string(colorCyan), l1, string(colorReset), l2, string(colorCyan), l3, string(colorReset))

	settings := []string{
		"fs.file-max = 2000000", "kernel.shmax = 60000000", "vm.overcommit_memory = 1",
	}

	filePath := "/etc/sysctl.d/99-illumio.conf"
	output := grepFiles(settings, filePath)
	createTable(output)
}

func grepFiles(settings []string, path string) []checkResult {
	list := make([]checkResult, 0)

	for _, raw := range settings {
		v := formatString(raw)
		output, err := exec.Command("/usr/bin/grep", v, path).Output()

		commandOutput := strings.TrimSuffix(string(output), "\n")
		configured := true

		if err != nil {
			commandOutput = "NIL"
			configured = false
		}

		list = append(list, checkResult{
			Filename: path, Name: raw, Value: commandOutput, Configured: configured,
		})
	}
	return list
}

// format the string
func formatString(str string) string {
	var rString string

	if strings.Contains(str, " = ") {
		rString = strings.Replace(str, " = ", ".*=.*", -1)
	} else {
		rString = strings.Replace(str, "=", ".*", -1)
	}

	return rString
}

// create table view on console
func createTable(raw []checkResult) {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Illumio Specific Path"},
			{Align: simpletable.AlignCenter, Text: "Required"},
			{Align: simpletable.AlignCenter, Text: "Discovered"},
			{Align: simpletable.AlignCenter, Text: "Configured?"},
		},
	}

	for _, v := range raw {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: v.Filename},
			{Align: simpletable.AlignLeft, Text: v.Name},
			{Align: simpletable.AlignLeft, Text: color(v.Value)},
			{Align: simpletable.AlignCenter, Text: color(fmt.Sprintf("%t", v.Configured))},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.SetStyle(simpletable.StyleRounded)
	fmt.Println(table.String())
}

func color(s string) string {
	if s == "NIL" || s == "false" {
		return fmt.Sprintf("%s%s%s", colorRed, s, colorDefault)
	}

	return s
}
