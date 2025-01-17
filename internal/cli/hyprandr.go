package cli

import (
	"fmt"
	"strings"

	"github.com/Tigy01/hyprandr/internal/monitors"
	myerrors "github.com/Tigy01/hyprandr/internal/myErrors"
)

type monitor = monitors.Monitor
type monitorMap = map[string]*monitor

func Run(currentMonitors monitorMap, resolution, refreshRate int, monitorName, customRes string, toggle bool) {
	if _, found := currentMonitors[monitorName]; !found {
		err := fmt.Sprintln("Must provide monitor name")
		for k := range currentMonitors {
			err += fmt.Sprintln(k)
		}
		fmt.Printf("err: %v\n", err)
		return
	}

	if resolution != -1 {
		myerrors.Try(ChangeRes(currentMonitors, monitorName, resolution))
		return
	}

	if refreshRate != -1 {
		myerrors.Try(SetRefresh(currentMonitors, monitorName, refreshRate))
		return
	}

	if customRes != "none" {
		myerrors.Try(SetRes(currentMonitors, monitorName, customRes))
		return
	}

	if toggle {
		myerrors.Try(ToggleMonitor(currentMonitors, monitorName))
		return
	}

	fmt.Println("\nmust add an additional arg")
}

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

func ToggleMonitor(currentMonitors monitorMap, monitorName string) error {
	monitor := currentMonitors[monitorName]
	monitor.Disable = !monitor.Disable
	return RewriteConfig(currentMonitors)
}

func ChangeRes(currentMonitors monitorMap, monitorName string, resIndex int) error {
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

func SetRefresh(currentMonitors monitorMap, monitorName string, refresh int) error {
	selection, ok := currentMonitors[monitorName]

	if !ok {
		return fmt.Errorf("Invalid monitor name")
	}

	res := strings.Split(selection.CurrentRes, "@")[0]
	fmt.Println(res, refresh)

	refreshIndex := getRefreshRateIndex(*selection, res, refresh)
	if refreshIndex == -1 {
		return fmt.Errorf("Invalid Refresh Rate: %v\n\tMust be an integer", refresh)
	}

	selection.CurrentRes = res + fmt.Sprintf("@%v", selection.Modes[res][refreshIndex])

	return RewriteConfig(currentMonitors)
}

func SetRes(currentMonitors monitorMap, monitorName string, newRes string) error {
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

func PrintCurrentConfig(currentMonitors monitorMap) {
	for name, monitor := range currentMonitors {
		hOffset := monitor.HOffset
		vOffset := monitor.VOffset
		scale := monitor.Scale
		res := monitor.CurrentRes
		fmt.Println(fmt.Sprintf("monitor=%s, %s, %vx%v, %v", name, res, hOffset, vOffset, scale))
	}
}

func PrintResolutions(currentMonitors monitorMap) {
	for name, monitor := range currentMonitors {
		fmt.Println(name)
		for i, res := range monitor.Resolutions {
			fmt.Println(fmt.Sprintf("%v: %s", i, res))
			fmt.Printf("\tRefresh Rates: %v\n", monitor.Modes[res])
		}
		fmt.Println("")
	}
}

func getRefreshRateIndex(selection monitor, res string, refresh int) int {
	for i, rate := range selection.Modes[res] {
		if rate == refresh {
			return i
		}
	}
	return -1
}
