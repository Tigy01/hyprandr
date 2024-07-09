package cli

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/Tigy01/hyprandr/internal/monitors"
)

var restOfFile []string

// Returns the path of the displays.conf file within the user's filesystem
func GetConfigPath() (path string, err error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/hypr/displays.conf", configDir), nil
}

// Returns a map of names to monitor structs with a variety of information
// about them
func GetCurrentSettings() (map[string]*monitor, error) {
	restOfFile = []string{}

	avaliableMonitors, err := monitors.GetMonitors()
	if err != nil {
		return nil, err
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	displayFile, err := os.Open(configPath)
	defer displayFile.Close()
	if err != nil {
		CreateDefaultConfig()
		return GetCurrentSettings()
	}

	lines := []string{}

	newMonitors := make([]string, 0)
	newOffset := "0"

	for name := range avaliableMonitors {
		newMonitors = append(newMonitors, name)
	}

	scanner := bufio.NewScanner(displayFile)
	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), " ", "")
		lines = append(lines, line)

		if line == "" {
			continue
		}

		if cutLine, found := strings.CutPrefix(line, "monitor="); found == true {
			name, monitor := parseMonitorLine(cutLine)

			avaliableMonitor, ok := avaliableMonitors[name]
			if ok {
				avaliableMonitor.CurrentRes = monitor.CurrentRes
				avaliableMonitor.Scale = monitor.Scale
				avaliableMonitor.HOffset = monitor.HOffset
				avaliableMonitor.VOffset = monitor.VOffset
				avaliableMonitor.OtherOptions = monitor.OtherOptions

				offset, err := strconv.ParseInt(
					strings.Split(avaliableMonitor.CurrentRes, "x")[0],
					10,
					64,
				)

				if err != nil {
					return nil, err
				}

				scale, err := strconv.ParseFloat(
					avaliableMonitor.Scale,
					64,
				)
				if err != nil {
					return nil, err
				}

				newOffset = fmt.Sprintf("%v", math.Round(float64(offset)/scale))

				if newIndex := slices.Index(newMonitors, name); newIndex > -1 {
					newMonitors = slices.Delete(newMonitors, newIndex, newIndex+1)
				}
				continue
			}
			if _, found := strings.CutSuffix(line, "#Invalid"); !found {
				line += " #Invalid"
			}
		}
		restOfFile = append(restOfFile, line)
	}

	for _, name := range newMonitors {
		avaliableMonitor := avaliableMonitors[name]

		avaliableMonitor.CurrentRes = avaliableMonitor.Resolutions[0]
		avaliableMonitor.CurrentRes = fmt.Sprintf(
			"%v@%v",
			avaliableMonitor.CurrentRes,
			avaliableMonitor.Modes[avaliableMonitor.CurrentRes][0],
		)

		avaliableMonitor.HOffset = newOffset
		avaliableMonitor.VOffset = "0"

		avaliableMonitor.Scale = "1"
		avaliableMonitor.OtherOptions = ""

		currentOffset, err := strconv.ParseInt(newOffset, 10, 64)
		if err != nil {
			return nil, err
		}
		addition, err := strconv.ParseInt(
			strings.Split(avaliableMonitor.CurrentRes,
				"x",
			)[0], 10, 64,
		)
		if err != nil {
			return nil, err
		}

		newOffset = fmt.Sprintf("%v", currentOffset+addition)
	}

	return avaliableMonitors, nil
}

func RewriteConfig(currentMonitors map[string]*monitor) error {
	config, err := GetConfigPath()
	if err != nil {
		return err
	}

	os.Remove(config)
	file, err := os.Create(config)
	defer file.Close()
	if err != nil {
		return err
	}

	for name, monitor := range currentMonitors {
		line := fmt.Sprintf(
			"monitor=%s, %s, %sx%s, %s, %s\n",
			name,
			monitor.CurrentRes,
			monitor.HOffset,
			monitor.VOffset,
			monitor.Scale,
			monitor.OtherOptions,
		)
		file.WriteString(line)
	}

	for _, line := range restOfFile {
		file.WriteString(line + "\n")
	}

	return nil
}

// Parses a hyprland formatted monitor line
func parseMonitorLine(line string) (name string, newMonitor *monitor) {
	name, line, _ = strings.Cut(line, ",")
	resolution, line, _ := strings.Cut(line, ",")
	hoffset, line, _ := strings.Cut(line, "x")
	voffset, line, _ := strings.Cut(line, ",")
	scale, other, _ := strings.Cut(line, ",")
	return name, &monitor{
		Resolutions:  []string{},
		CurrentRes:   resolution,
		HOffset:      hoffset,
		VOffset:      voffset,
		Scale:        scale,
		OtherOptions: other,
	}
}
