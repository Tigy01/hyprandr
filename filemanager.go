package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getConfigPath() (path string, err error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/hypr/displays.conf", configDir), nil
}

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
	name := ""
	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), " ", "")
		lines = append(lines, line)

		if line == "" {
			continue
		}

		if strings.Contains(line, "monitor") {
			name, monitor := parseMonitorLine(line)
			monitors[name] = monitor
			continue
		}

		if delimiter := strings.Index(line, ":"); delimiter > 0 {
			name = line[delimiter+1:]
			continue
		}

		m, ok := monitors[name]
		if !ok {
			continue
		}
		m.resolutions = append(m.resolutions, line[1:])
	}

	return monitors, nil
}

func saveResolutions(monitors map[string]*monitor) {
	configPath, err := getConfigPath()
	file, err2 := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil || err2 != nil {
		fmt.Println("could not write res")
		return
	}
	defer file.Close()

	file.WriteString("# Current Resolutions #\n")
	for k, m := range monitors {
		resolutions := m.sortResolutions()

		_, err = file.WriteString("\n#Name:" + fmt.Sprint(k) + "\n")
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		for _, r := range resolutions {
			file.WriteString("#" + fmt.Sprintln(r))
		}
	}
}

func rewriteConfig(currentMonitors map[string]*monitor) {
	config, err := getConfigPath()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	os.Remove(config)
	file, err := os.Create(config)
	defer file.Close()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	for name, monitor := range currentMonitors {
		line := fmt.Sprintf("monitor=%s, %s, %sx%s, %s\n", name, monitor.currentRes, monitor.hOffset, monitor.vOffset, monitor.scale)
		file.WriteString(line)
	}
	saveResolutions(currentMonitors)
}

func parseMonitorLine(line string) (name string, newMonitor *monitor) {
	line = line[strings.Index(line, "=")+1:]
	name = line[:strings.Index(line, ",")]

	line = line[strings.Index(line, ",")+1:]
	resolution := line[:strings.Index(line, ",")]

	line = line[strings.Index(line, ",")+1:]
	hoffset := line[:strings.Index(line, "x")]

	line = line[strings.Index(line, "x")+1:]
	voffset := line[:strings.Index(line, ",")]

	line = line[strings.Index(line, ",")+1:]
	scale := line
	return name, &monitor{
		resolutions: []string{},
		currentRes:  resolution,
		hOffset:     hoffset,
		vOffset:     voffset,
		scale:       scale,
	}
}
