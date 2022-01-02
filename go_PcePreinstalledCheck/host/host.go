package host

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/alexeyco/simpletable"
)

type hostInfo struct {
	Desc     string
	Current  string
	Required string
	Set      bool
}

const (
	colorCyan    = "\033[36m"
	colorReset   = "\033[0m"
	colorDefault = "\x1b[39m"
	colorRed     = "\x1b[91m"
)

func HostCheck(cpuNeeded, ramNeeded int) {
	list := make([]hostInfo, 0)
	set := true
	fmt.Println(string(colorCyan), "\n\nBasic Host Check:", string(colorReset))
	// hostname
	cmd, req := "hostname", "-"
	output, _ := exec.Command(cmd).Output()
	commandOutput("hostname", output, req, &list, set)

	// ip adderss
	output, _ = exec.Command(cmd, "-I").Output()
	commandOutput("IP Address", output, req, &list, set)

	// os version
	req = "7.4 above | 8.3 above"
	output, _ = exec.Command("bash", "-c", "cat /etc/redhat-release | rev | cut -d'(' -f 2 | rev | awk 'NF>1{print $NF}'").Output()
	parse := strings.TrimSuffix(string(output), "\n")
	f, _ := strconv.ParseFloat(parse, 32)

	if f < 7.4 {
		set = false
	}
	commandOutput("Supported Version", output, req, &list, set)

	// locale set
	req = "en_US.UTF-8 | en_GB.UTF-8"
	output, _ = exec.Command("bash", "-c", "grep -r -E -i '(en_US.UTF-8|en_GB.UTF-8)' /etc/locale.conf").Output()
	parse = strings.TrimSuffix(string(output), "\n")

	if !(strings.Contains(parse, "en_US") || strings.Contains(parse, "en_GB")) {
		set = false
	}
	commandOutput("Locale", output, req, &list, set)

	// cpu check
	req = strconv.Itoa(cpuNeeded)
	output, _ = exec.Command("bash", "-c", "getconf _NPROCESSORS_ONLN").Output()
	parse = strings.TrimSuffix(string(output), "\n")
	i, _ := strconv.Atoi(parse)

	if i < cpuNeeded {
		set = false
	}
	commandOutput("Number of CPUs", output, req, &list, set)

	// ram check
	req = strconv.Itoa(ramNeeded)
	output, _ = exec.Command("bash", "-c", "awk '$1 == \"MemTotal:\" { print $2 / 1000 / 1000 }' /proc/meminfo").Output()
	parse = strings.TrimSuffix(string(output), "\n")
	i, _ = strconv.Atoi(parse)

	if i < ramNeeded {
		set = false
	}
	commandOutput("RAM per Node (GB)", output, req, &list, set)

	createTable(list)

}

func commandOutput(desc string, input []byte, req string, list *[]hostInfo, set bool) {
	commandOutput := strings.TrimSuffix(string(input), "\n")

	*list = append(*list, hostInfo{
		Desc: desc, Current: commandOutput, Required: req, Set: set,
	})

}

// create table view on console
func createTable(raw []hostInfo) {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Desc"},
			{Align: simpletable.AlignCenter, Text: "Current"},
			{Align: simpletable.AlignCenter, Text: "Required"},
			{Align: simpletable.AlignCenter, Text: "Set?"},
		},
	}

	for _, v := range raw {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: v.Desc},
			{Align: simpletable.AlignLeft, Text: v.Current},
			{Align: simpletable.AlignLeft, Text: v.Required},
			{Align: simpletable.AlignCenter, Text: color(fmt.Sprintf("%t", v.Set))},
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
