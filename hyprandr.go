package main

import (
	"fmt"
	"strings"
)

func createDefaultConfig() error {
	monitors, err := getMonitors()

	if err != nil {
		return err
	}

	defaultHOffset := "0"
	for _, monitor := range monitors {
		monitor.currentRes = monitor.resolutions[0]
		monitor.currentRes += fmt.Sprintf("@%v", monitor.modes[monitor.currentRes][0])
		monitor.hOffset = defaultHOffset
		monitor.vOffset = "0"
		monitor.scale = "1"

		defaultHOffset = monitor.currentRes[:strings.Index(monitor.currentRes, "x")]
	}
	return rewriteConfig(monitors)
}

func setRes(currentMonitors map[string]*monitor, monitorName string, newRes string) error {
	selection, ok := currentMonitors[monitorName]

	if strings.Index(newRes, "@") == -1 {
		return fmt.Errorf(
			"Must include a refresh rate in the form of @[REFRESH RATE]",
		)
	}

	if !ok {
		return fmt.Errorf("Invalid Monitor Name: %v", monitorName)
	}

	selection.currentRes = newRes
	return rewriteConfig(currentMonitors)
}

func changeRes(currentMonitors map[string]*monitor, monitorName string, resIndex int) error {
	selection, ok := currentMonitors[monitorName]

	if !ok {
		return fmt.Errorf("Invalid monitor name")
	}

	for i, res := range currentMonitors[monitorName].resolutions {
		if i == resIndex {
			selection.currentRes = res + fmt.Sprintf("@%v", selection.modes[res][0])
		}
	}

	return rewriteConfig(currentMonitors)
}

func printCurrentConfig(currentMonitors map[string]*monitor) {
	for name, monitor := range currentMonitors {
		hOffset := monitor.hOffset
		vOffset := monitor.vOffset
		scale := monitor.scale
		res := monitor.currentRes
		fmt.Println(fmt.Sprintf("monitor=%s, %s, %vx%v, %v", name, res, hOffset, vOffset, scale))
	}
}

func printResolutions(monitors map[string]*monitor) {
	for name, monitor := range monitors {
		fmt.Println(name)
		for i, res := range monitor.resolutions {
			fmt.Println(fmt.Sprintf("%v: %s", i, res))
		}
		fmt.Println("")
	}
}
