package menu

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/manifoldco/promptui"
)

type PCE struct {
	Description string
	CpuPerNode  int
	RamPerNode  int
	Notes       string
}

func PceTypeMenu() (int, int, int) {
	// PCE Spec init
	// 0 - SNC, 1 - 2x2 small, 2 - 2x2 regular, 3 - 4x2
	pceVirtualSpec := []PCE{
		{
			Description: "SNC (Single Node Cluster)", CpuPerNode: 6, RamPerNode: 16,
			Notes: "\n3 CPU cores per node IF running on physical machine\n- Intel® Xeon(R) CPU E5-2695 v4 at 2.10GHz or equivalent",
		},
		{
			Description: "2x2 MNC (Multi Node Cluster for < 2500 VENs)", CpuPerNode: 8, RamPerNode: 32,
			Notes: "\n4 CPU cores per node IF running on physical machine\n- Intel® Xeon(R) CPU E5-2695 v4 at 2.10GHz or equivalent",
		},
		{
			Description: "2x2 MNC (Multi Node Cluster for < 10000 VENs)", CpuPerNode: 32, RamPerNode: 64,
			Notes: "\n128GB RAM Recommended.\n16 CPU cores per node IF running on physical machine\n- Intel® Xeon(R) CPU E5-2695 v4 at 2.10GHz or equivalent",
		},
		{
			Description: "4x2 MNC (High Spec Multi Node Cluster for < 25000 VENs)", CpuPerNode: 32, RamPerNode: 128,
			Notes: "\n16 CPU cores per node IF running on physical machine\n- Intel® Xeon(R) CPU E5-2695 v4 at 2.10GHz or equivalent",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}:",
		Active:   "→ {{ .Description | green }}",
		Inactive: "   {{ .Description | faint }}",
		Selected: "\U0001F44D {{ .Description | red | cyan }}",
		Help:     "Use the arrow keys to navigate: ↓ ↑",
		Details: `

--- PCE Spec on VM -------------------------------------------------------------------
{{ "Description:" | cyan }}	{{ .Description }}
{{ "Num of vCPU per node:" | cyan }}	{{ .CpuPerNode }} 
{{ "Num of RAM(GB) per node:" | cyan }}	{{ .RamPerNode }}
{{ "Notes:" | cyan }}	{{ .Notes | yellow }}`,
	}

	prompt := promptui.Select{
		Label:     "Select PCE Model",
		Items:     pceVirtualSpec,
		Templates: templates,
	}

	// clear the screen
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("You choose PCE cluster type %d: %s\n", i+1, pceVirtualSpec[i].Description)
	return i, pceVirtualSpec[i].CpuPerNode, pceVirtualSpec[i].RamPerNode
}
