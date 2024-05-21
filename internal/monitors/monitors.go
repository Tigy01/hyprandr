package monitors 

import (
	"math"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

type Monitor struct {
	Resolutions  []string
	Modes        map[string][]int //resolution to refresh rates
	CurrentRes   string
	HOffset      string
	VOffset      string
	Scale        string
	OtherOptions string
}

// Gets the Monitor Modes from the 'hyprctl monitors all' command
//
// returns map of names to monitors containing lists of Resolutions and
// map of Resolutions to Modes
func GetMonitors() (map[string]*Monitor, error) {
	rawSystemInfo, err := exec.Command("hyprctl", "monitors", "all").Output()
	if err != nil {
		return nil, err
	}

	systemInfo := strings.Split(string(rawSystemInfo), "\n")
	monitors := map[string]*Monitor{}

	var currentMonitor *Monitor
	for _, line := range systemInfo {
		monitorIndex := strings.Index(line, "Monitor ")

		//arbitrary range meant to eliminate edge cases
		if monitorIndex > -1 && monitorIndex < 2 {
			line = line[monitorIndex+8:]
			name := line[:strings.Index(line, " ")]
			monitors[name] = &Monitor{
				Modes: make(map[string][]int, 0),
			}
			currentMonitor = monitors[name]
			continue
		}

		parseModes(line, currentMonitor)
		for res := range currentMonitor.Modes {
			slices.Sort(currentMonitor.Modes[res])
			slices.Reverse(currentMonitor.Modes[res])
		}
	}
	return monitors, nil
}

func parseModes(line string, currentMonitor *Monitor) {
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
		if !slices.Contains(currentMonitor.Resolutions, resolution) {
			currentMonitor.Resolutions = append(currentMonitor.Resolutions, resolution)
		}
		currentMonitor.Modes[resolution] = append(currentMonitor.Modes[resolution], finalRate)
	}
}
