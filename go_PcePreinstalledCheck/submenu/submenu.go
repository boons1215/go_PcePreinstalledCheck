package submenu

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

type nodeType struct {
	Type string
}

func NodeTypeMenu() int {
	// node type init
	nodeType := []nodeType{
		{Type: "Core"},
		{Type: "Data"},
	}

	templates := &promptui.SelectTemplates{
		Label:    "Select Node Type:",
		Active:   "→ {{ .Type | green }}",
		Inactive: "   {{ .Type | faint }}",
		Selected: "\U0001F44D {{ .Type | red | cyan }}"}

	prompt := promptui.Select{
		Items:     nodeType,
		Templates: templates,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("You choose node type %d: %s\n", i+1, nodeType[i].Type)

	return i
}
