package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Returns the path of the displays.conf file within the user's filesystem
func getConfigPath() (path string, err error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/hypr/displays.conf", configDir), nil
}

// Returns a map of names to monitor structs with a variety of information
// about them
func getCurrentSettings() (map[string]*monitor, error) {
	monitors := map[string]*monitor{}

	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	displayFile, err := os.Open(configPath)
	defer displayFile.Close()
	if err != nil {
		createDefaultConfig()
		return getCurrentSettings()
	}

	lines := []string{}

	scanner := bufio.NewScanner(displayFile)
	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), " ", "")
		lines = append(lines, line)

		if line == "" {
			continue
		}

		if line, found := strings.CutPrefix(line, "monitor="); found == true {
			name, monitor := parseMonitorLine(line)
			monitors[name] = monitor
			continue
		}
	}

	avaliableMonitors, err := getMonitors()
	for name, monitor := range monitors {
		if err != nil {
			return nil, err
		}
		monitor.modes = avaliableMonitors[name].modes
		monitor.resolutions = avaliableMonitors[name].resolutions
	}

	return monitors, nil
}

func rewriteConfig(currentMonitors map[string]*monitor) error {
	config, err := getConfigPath()
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
		line := fmt.Sprintf("monitor=%s, %s, %sx%s, %s\n", name, monitor.currentRes, monitor.hOffset, monitor.vOffset, monitor.scale)
		file.WriteString(line)
	}
	return nil
}

// Parses a hyprland formatted monitor line
func parseMonitorLine(line string) (name string, newMonitor *monitor) {
	name, line, _ = strings.Cut(line, ",")
	resolution, line, _ := strings.Cut(line, ",")
	hoffset, line, _ := strings.Cut(line, "x")
	voffset, line, _ := strings.Cut(line, ",")
	scale, _, _ := strings.Cut(line, ",")
	return name, &monitor{
		resolutions: []string{},
		currentRes:  resolution,
		hOffset:     hoffset,
		vOffset:     voffset,
		scale:       scale,
	}
}
