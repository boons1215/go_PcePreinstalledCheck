package pkgcheck

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/alexeyco/simpletable"
)

type checkResult struct {
	Name      string
	Value     string
	Installed bool
}

const (
	colorCyan    = "\033[36m"
	colorReset   = "\033[0m"
	colorDefault = "\x1b[39m"
	colorRed     = "\x1b[91m"
)

func Execute() {
	go rpmPkgCheck()
	fipsPkgCheck()
	libPkgCheck()
}

func rpmPkgCheck() {
	packages := []string{
		"bash", "bzip2", "chkconfig", "coreutils", "findutils", "curl", "gawk", "grep",
		"initscripts", "gzip", "logrotate", "net-tools", "procps-ng", "sed", "shadow-utils",
		"tar", "util-linux", "zlib", "bind-utils", "libnfnetlink", "glibc", "ncurses",
	}

	rpmList := multiThreadRpmCheck(packages)
	fmt.Println(string(colorCyan), "\n\nRequired packages/dependencies for PCE:", string(colorReset))

	// create table
	createTable(rpmList)
}

func libPkgCheck() {
	packages := []string{
		"glibc", "libgcc", "libstdc++", "libuuid", "ncurses-libs", "openssl", "zlib",
	}

	files := []string{
		"libreadline.so.6", "libselinux.so.1",
	}

	libList := multiThreadRpmCheck(packages)
	fmt.Println(string(colorCyan), "\n\nRequired Libraries for PCE:", string(colorReset))

	fileList := findFiles(files)
	libList = append(libList, fileList...)
	// create table
	createTable(libList)
}

func fipsPkgCheck() {
	packages := []string{
		"libssl.so.10", "libcrypto.so.10",
	}

	fipsList := findFiles(packages)
	fmt.Println(string(colorCyan), "\n\nRequired FIPS Package for PCE:", string(colorReset))

	// create table
	createTable(fipsList)
}

// find files in /usr/lib64/
func findFiles(files []string) []checkResult {
	list := make([]checkResult, 0)

	for _, v := range files {
		output, err := exec.Command("/usr/bin/find", "/usr/lib64/", "-iname", v).Output()

		commandOutput := strings.TrimSuffix(string(output), "\n")
		installed := true

		// if package not installed
		if err != nil {
			commandOutput = "NIL"
			installed = false
		}

		list = append(list, checkResult{
			Name: v, Value: commandOutput, Installed: installed,
		})
	}
	return list
}

// multi thread rpm check function
func multiThreadRpmCheck(packages []string) []checkResult {
	list := make([]checkResult, 0)

	// multi-thread
	var wg sync.WaitGroup
	wg.Add(len(packages))

	for i := 0; i < len(packages); i++ {
		go func(i int) {
			defer wg.Done()
			output, err := exec.Command("/bin/rpm", "-q", packages[i]).Output()

			commandOutput := strings.TrimSuffix(string(output), "\n")
			installed := true

			// if package not installed
			if err != nil {
				commandOutput = "NIL"
				installed = false
			}

			list = append(list, checkResult{
				Name: packages[i], Value: commandOutput, Installed: installed,
			})
		}(i)
	}

	wg.Wait()
	return list
}

// create table view on console
func createTable(raw []checkResult) {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Required Package"},
			{Align: simpletable.AlignCenter, Text: "Discovered"},
			{Align: simpletable.AlignCenter, Text: "Installed?"},
		},
	}

	for _, v := range raw {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: v.Name},
			{Text: color(v.Value)},
			{Align: simpletable.AlignCenter, Text: color(fmt.Sprintf("%t", v.Installed))},
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
