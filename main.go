package main

import (
	"log"
	"os/exec"
	"strconv"

	"github.com/boons1215/goPCESpecChk/host"
	"github.com/boons1215/goPCESpecChk/kernel"
	"github.com/boons1215/goPCESpecChk/menu"
	"github.com/boons1215/goPCESpecChk/pkgcheck"
	"github.com/boons1215/goPCESpecChk/submenu"
)

func main() {
	// only root allowed to execute
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	checkErr(err)

	// 0 = root
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	checkErr(err)

	if i == 0 {
		// call main menu
		// 0 - SNC, 1 - 2x2 small, 2 - 2x2 regular, 3 - 4x2
		clusterType, cpuNeeded, ramNeeded := menu.PceTypeMenu()

		// call sub menu
		// 0 for Core, 1 for Data, 2 for SNC (default)
		nodeType := 2
		if clusterType != 0 {
			nodeType = submenu.NodeTypeMenu()
		}

		host.HostCheck(cpuNeeded, ramNeeded)
		pkgcheck.Execute()
		kernel.Execute(clusterType, nodeType)
	} else {
		log.Fatal("This script must be run as root or sudo!")
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
