package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Tigy01/hyprandr/internal/cli"
	"github.com/Tigy01/hyprandr/internal/monitors"
	"github.com/Tigy01/hyprandr/internal/myErrors"
	"github.com/Tigy01/hyprandr/internal/pages"
	"github.com/charmbracelet/bubbletea"
)

var windowWidth int
var windowHeight int
var currentMonitors map[string]*monitors.Monitor

func main() {
	currentMonitors = myerrors.TryWithValue(cli.GetCurrentSettings())

	var selection, refreshRate int
	var monitorName, customRes string
	var toggle bool
	if parseFlags(&selection, &refreshRate, &monitorName, &customRes, &toggle) {
		cli.Run(currentMonitors, selection, refreshRate, monitorName, customRes, toggle)
		return
	}

	if len(os.Args[1:]) != 0 {
		fmt.Printf(
			"Error: Unrecognised flags\n\tMust use the -cli flag for cli commands: '%v'\n",
			os.Args[1:],
		)
		return
	}

	tea.NewProgram(
		pages.MonitorSelectPage{}.New(currentMonitors),
		tea.WithAltScreen(),
	).Run()

}

// Returns True if the program is running in cliMode
func parseFlags(selection, refreshRate *int, monitorName, customRes *string, toggle *bool) (cliMode bool) {
	flag.BoolFunc("cli", "used to enter cli mode", func(s string) error { cliMode = true; return nil })
	flag.IntVar(selection, "change-res", -1, "select the cooresponding resolution")
	flag.StringVar(monitorName, "monitor", "none", "used to select the monitor to change")
	flag.StringVar(customRes, "set-res", "none", "used to set the current resolution")
	flag.IntVar(refreshRate, "change-refresh", -1, "used to set the current refresh rate")
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
            os.Exit(0)
			return nil
		},
	)
	flag.Parse()
	return
}
