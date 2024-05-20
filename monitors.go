package main

import (
	"math"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

type monitor struct {
	resolutions  []string
	modes        map[string][]int //resolution to refresh rates
	currentRes   string
	hOffset      string
	vOffset      string
	scale        string
	otherOptions string
}

// Gets the monitor modes from the 'hyprctl monitors all' command
//
// returns map of names to monitors containing lists of resolutions and
// map of resolutions to modes
func getMonitors() (map[string]*monitor, error) {
	rawSystemInfo, err := exec.Command("hyprctl", "monitors", "all").Output()
	if err != nil {
		return nil, err
	}

	systemInfo := strings.Split(string(rawSystemInfo), "\n")
	monitors := map[string]*monitor{}

	var currentMonitor *monitor
	for _, line := range systemInfo {
		monitorIndex := strings.Index(line, "Monitor ")

		//arbitrary range meant to eliminate edge cases
		if monitorIndex > -1 && monitorIndex < 2 {
			line = line[monitorIndex+8:]
			name := line[:strings.Index(line, " ")]
			monitors[name] = &monitor{
				modes: make(map[string][]int, 0),
			}
			currentMonitor = monitors[name]
			continue
		}

		parseModes(line, currentMonitor)
		for res := range currentMonitor.modes {
			slices.Sort(currentMonitor.modes[res])
			slices.Reverse(currentMonitor.modes[res])
		}
	}
	return monitors, nil
}

func parseModes(line string, currentMonitor *monitor) {
	modeIndex := strings.Index(line, "availableModes: ")

	if modeIndex == -1 {
		return
	}

	modeList := strings.Split(line[modeIndex+16:], " ")
	for _, element := range modeList {
		resolution, rate, _ := strings.Cut(element, "@")
		rate, _ = strings.CutSuffix(rate, "Hz")
		convertedRate, err := strconv.ParseFloat(rate, 32)
		if err != nil {
			continue
		}
		finalRate := int(math.Round(convertedRate))
		if !slices.Contains(currentMonitor.resolutions, resolution) {
			currentMonitor.resolutions = append(currentMonitor.resolutions, resolution)
		}
		currentMonitor.modes[resolution] = append(currentMonitor.modes[resolution], finalRate)
	}
}
