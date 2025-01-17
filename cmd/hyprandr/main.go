package main

import (
	"flag"
	"fmt"
	"github.com/Tigy01/hyprandr/internal/cli"
	"github.com/Tigy01/hyprandr/internal/pages"
	tea "github.com/charmbracelet/bubbletea"
	"os"
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
	var toggle bool
	cmdMode := false
	flag.BoolFunc("cmd", "used to enter cmd mode", func(s string) error { cmdMode = true; return nil })
	flag.IntVar(&selection, "change-res", -1, "select the cooresponding resolution")
	flag.StringVar(&monitorName, "monitor", "none", "used to select the monitor to change")
	flag.StringVar(&customRes, "set-res", "none", "used to set the current resolution")
	flag.BoolFunc("toggle",
		"used to toggle a monitor from the command line",
		func(s string) error { toggle = true; return nil },
	)

	flag.BoolFunc(
		"list",
		"Used to list all current resolutions",
		func(s string) error {
			cli.PrintResolutions(currentMonitors)
			return nil
		},
	)
	flag.Parse()

	if !cmdMode {
		if len(os.Args[1:]) != 0 {
			fmt.Printf("Error: Unrecognised flags: %v\n", os.Args[1:])
            return
		}

		tea.NewProgram(
			pages.MonitorSelectPage{}.New(currentMonitors), tea.WithAltScreen(),
		).Run()
		return

	}

	if _, found := currentMonitors[monitorName]; !found {
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

	if toggle {
		err = cli.ToggleMonitor(currentMonitors, monitorName)

		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		return
	}

	fmt.Println("\nmust add an additional arg")
}
