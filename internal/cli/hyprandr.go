package cli

import (
	"fmt"
	"github.com/Tigy01/hyprandr/internal/monitors"
	"strings"
)

type monitor = monitors.Monitor

func CreateDefaultConfig() error {
	currentMonitors, err := monitors.GetMonitors()

	if err != nil {
		return err
	}

	defaultHOffset := "0"
	for _, monitor := range currentMonitors {
		monitor.CurrentRes = monitor.Resolutions[0]
		monitor.CurrentRes += fmt.Sprintf("@%v", monitor.Modes[monitor.CurrentRes][0])
		monitor.HOffset = defaultHOffset
		monitor.VOffset = "0"
		monitor.Scale = "1"

		defaultHOffset = monitor.CurrentRes[:strings.Index(monitor.CurrentRes, "x")]
	}
	return RewriteConfig(currentMonitors)
}

func disableMonitor() error {
    return nil
}

func ChangeRes(currentMonitors map[string]*monitor, monitorName string, resIndex int) error {
	selection, ok := currentMonitors[monitorName]

	if !ok {
		return fmt.Errorf("Invalid monitor name")
	}

	for i, res := range currentMonitors[monitorName].Resolutions {
		if i == resIndex {
			selection.CurrentRes = res + fmt.Sprintf("@%v", selection.Modes[res][0])
		}
	}

	return RewriteConfig(currentMonitors)
}

func SetRes(currentMonitors map[string]*monitor, monitorName string, newRes string) error {
	selection, ok := currentMonitors[monitorName]
	if strings.Index(newRes, "@") == -1 {
		return fmt.Errorf(
			"Must include a refresh rate in the form of @[REFRESH RATE]",
		)
	}

	if !ok {
		return fmt.Errorf("Invalid Monitor Name: %v", monitorName)
	}

	selection.CurrentRes = newRes
	return RewriteConfig(currentMonitors)
}

func PrintCurrentConfig(currentMonitors map[string]*monitor) {
	for name, monitor := range currentMonitors {
		hOffset := monitor.HOffset
		vOffset := monitor.VOffset
		scale := monitor.Scale
		res := monitor.CurrentRes
		fmt.Println(fmt.Sprintf("monitor=%s, %s, %vx%v, %v", name, res, hOffset, vOffset, scale))
	}
}

func PrintResolutions(currentMonitors map[string]*monitor) {
	for name, monitor := range currentMonitors {
		fmt.Println(name)
		for i, res := range monitor.Resolutions {
			fmt.Println(fmt.Sprintf("%v: %s", i, res))
		}
		fmt.Println("")
	}
}
