package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Tigy01/hyprandr/internal/cli"
	"github.com/Tigy01/hyprandr/internal/monitors"
	myerrors "github.com/Tigy01/hyprandr/internal/myErrors"
	"github.com/Tigy01/hyprandr/internal/pages"
	tea "github.com/charmbracelet/bubbletea"
)

var windowWidth int
var windowHeight int
var currentMonitors map[string]*monitors.Monitor

func main() {
	currentMonitors = myerrors.TryWithValue(cli.GetCurrentSettings())

	var selection int
	var monitorName, customRes string
	var toggle bool
	if parseFlags(&selection, &monitorName, &customRes, &toggle) {
		cli.Run(currentMonitors, selection, monitorName, customRes, toggle)
		return
	}

	if len(os.Args[1:]) != 0 {
		fmt.Printf(
			"Error: Unrecognised flags\n\tmust use the -cli flag for cli commands: %v\n",
			os.Args[1:],
		)
	}

	tea.NewProgram(
		pages.MonitorSelectPage{}.New(currentMonitors),
		tea.WithAltScreen(),
	).Run()

}

// Returns True if the program is running in cmdMode
func parseFlags(selection *int, monitorName, customRes *string, toggle *bool) (cmdMode bool) {
	flag.BoolFunc("cmd", "used to enter cmd mode", func(s string) error { cmdMode = true; return nil })
	flag.IntVar(selection, "change-res", -1, "select the cooresponding resolution")
	flag.StringVar(monitorName, "monitor", "none", "used to select the monitor to change")
	flag.StringVar(customRes, "set-res", "none", "used to set the current resolution")
	flag.BoolFunc("toggle",
		"used to toggle a monitor from the command line",
		func(s string) error {
			*toggle = true
			return nil
		},
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
	return
}
