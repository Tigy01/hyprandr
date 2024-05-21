package main

import (
	"flag"
	"fmt"
	"os"
    "github.com/Tigy01/hyprandr/internal/pages"
    "github.com/Tigy01/hyprandr/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
)

var windowWidth int
var windowHeight int

func main() {
	currentMonitors, err := cli.GetCurrentSettings()
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
			cli.PrintResolutions(currentMonitors)
			return nil
		},
	)
	flag.Parse()

	if !cmdMode && len(os.Args[1:]) == 0 {
		tea.NewProgram(
			pages.MonitorSelectPage{}.New(currentMonitors), tea.WithAltScreen(),
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
		err := cli.ChangeRes(currentMonitors, monitorName, selection)

		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		return
	}

	if customRes != "none" {
		err := cli.SetRes(currentMonitors, monitorName, customRes)

		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		return
	}

	fmt.Println("\nmust add an additional arg")
	return
}
