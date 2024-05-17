package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var windowWidth int
var windowHeight int

func main() {
	currentMonitors, err := getCurrentSettings()
	if err != nil {
		return
	}

	var selection int
	var monitorName string
	var customRes string
	cmdMode := false
	flag.BoolFunc("cmd", "used to enter cmd mode", func(s string) error { cmdMode = true; return nil })
	flag.IntVar(&selection, "change-res", -1, "select the cooresponding resolution")
	flag.StringVar(&monitorName, "monitor", "none", "used to select the monitor to change")
	flag.StringVar(&customRes, "set-res", "none", "used to set the current resolution")

	flag.BoolFunc(
		"list",
		"Used to list all current resolutions",
		func(s string) error {
			printResolutions(currentMonitors)
			return nil
		},
	)
	flag.Parse()

	if !cmdMode && len(os.Args[1:]) == 0 {
		tea.NewProgram(
			monitorSelectPage{}.New(currentMonitors), tea.WithAltScreen(),
		).Run()
		return
	}

	if monitorName == "none" {
		err := fmt.Sprintln("Must provide monitor name")
		for k := range currentMonitors {
			err += fmt.Sprintln(k)
		}
		fmt.Printf("err: %v\n", err)
		return
	}

	if selection != -1 {
		changeRes(currentMonitors, monitorName, selection)
		return
	}

	if customRes != "none" {
		setRes(currentMonitors, monitorName, customRes)
		return
	}

	fmt.Println("\nmust add an additional arg")
	return

}

func getHelp() string {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
        Render(" HELP: H -> Left, J -> Down, K -> Up, L -> Right, RETURN -> Select, Q -> Back ")
}
