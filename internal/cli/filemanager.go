package cli

import (
	"bufio"
	"fmt"
	"os"
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
	currentMonitors := map[string]*monitor{}

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

	scanner := bufio.NewScanner(displayFile)
	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), " ", "")
		lines = append(lines, line)

		if line == "" {
			continue
		}

		if cutLine, found := strings.CutPrefix(line, "monitor="); found == true {
			name, monitor := parseMonitorLine(cutLine)

			if avaliableMonitor, ok := avaliableMonitors[name]; ok {
				currentMonitors[name] = monitor
				currentMonitors[name].Resolutions = avaliableMonitor.Resolutions
				currentMonitors[name].Modes = avaliableMonitor.Modes
				continue
			}

			if _, found := strings.CutSuffix(line, "#Invalid"); !found {
				line += " #Invalid"
			}

			fmt.Println("Error: Invalid Monitor name in config:", name)
		}
		restOfFile = append(restOfFile, line)
	}

	return currentMonitors, nil
}

func rewriteConfig(currentMonitors map[string]*monitor) error {
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
