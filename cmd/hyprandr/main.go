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
		fmt.Printf("err: %v\n", err)
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
		err := changeRes(currentMonitors, monitorName, selection)

		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		return
	}

	if customRes != "none" {
		err := setRes(currentMonitors, monitorName, customRes)

		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		return
	}

	fmt.Println("\nmust add an additional arg")
	return

}

func getHelp() string {
	return lipgloss.JoinHorizontal(lipgloss.Bottom,
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true, false, true, true).
			Foreground(lipgloss.Color("#5555ff")).
			Render(" HELP "+lipgloss.NormalBorder().Right, " "),
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true, true, true, false).
			Render("H -> Left | J -> Down | K -> Up | L -> Right | RETURN -> Select | Q -> Back "),
	)
}
