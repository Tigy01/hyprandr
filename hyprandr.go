package main

import (
	"fmt"
	"strings"
)

func setRes(currentMonitors map[string]*monitor, monitorName string, newRes string) {
	selection, ok := currentMonitors[monitorName]

	if strings.Index(newRes, "@") == -1 {
		fmt.Println("Must include a refresh rate in the form of @[REFRESH RATE]")
		return
	}

	if !ok {
		fmt.Println("Invalid monitor name")
		return
	}

	selection.currentRes = newRes
	rewriteConfig(currentMonitors)
}

func changeRes(currentMonitors map[string]*monitor, monitorName string, resIndex int) {
	selection, ok := currentMonitors[monitorName]

	if !ok {
		fmt.Println("Invalid monitor name")
		return
	}

	for i, res := range currentMonitors[monitorName].resolutions {
		if i == resIndex {
			selection.currentRes = res + fmt.Sprintf("@%v", selection.modes[res][0]) 
		}
	}
    
	rewriteConfig(currentMonitors)
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

func createDefaultConfig() {
	monitors, err := getMonitors()

	if err != nil {
		return
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
	rewriteConfig(monitors)
}
