package main

import (
	"fmt"
	"math"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

type monitor struct {
	resolutions []string
	currentRes  string
	hOffset     string
	vOffset     string
	scale       string
}

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
			selection.currentRes = res
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
	monitorInfo, _ := exec.Command("xrandr").Output()
	monitorList := strings.Split(string(monitorInfo), "\n")
	monitorList = monitorList[1:]

	monitors := map[string]*monitor{}
	var currentmonitor string
	for _, line := range monitorList {
		if len(line) <= 1 {
			continue
		}

		if nameIndex := strings.Index(line, "connected"); nameIndex > 0 {
			currentmonitor = strings.ReplaceAll(line[:nameIndex-1], " ", "")
			monitors[currentmonitor] = &monitor{}
			continue
		}

		resLine, err := parseResLine(line)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		monitors[currentmonitor].resolutions = append(monitors[currentmonitor].resolutions, resLine)
	}
    
    defaultHOffset:="0"
	for _, monitor := range monitors {
		monitor.sortResolutions()
		monitor.currentRes = monitor.resolutions[0]
		monitor.hOffset = defaultHOffset 
		monitor.vOffset = "0"
		monitor.scale = "1"

        defaultHOffset = monitor.currentRes[:strings.Index(monitor.currentRes, "x")]
	}
	rewriteConfig(monitors)
}

func getResolutions(monitors map[string]*monitor) map[string]*monitor {
	monitorInfo, _ := exec.Command("xrandr").Output()
	monitorList := strings.Split(string(monitorInfo), "\n")
	monitorList = monitorList[1:]

	var currentmonitor string

	for _, v := range monitorList {
		if len(v) <= 1 {
			continue
		}

		if spaceIndex := strings.Index(v, " "); spaceIndex > 0 {
			currentmonitor = v[:spaceIndex]
			monitors[currentmonitor].resolutions = []string{}
			continue
		}

		resLine, err := parseResLine(v)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return nil
		}

		monitors[currentmonitor].resolutions = append(monitors[currentmonitor].resolutions, resLine)
	}

	for _, v := range monitors {
		v.sortResolutions()
	}

	return monitors
}

func parseResLine(line string) (string, error) {
	hres := strings.ReplaceAll(
		line[:strings.Index(line, "x")],
		" ",
		"",
	) //"   1920x1080   75   60" -> "1920"

	line = line[strings.Index(line, "x")+1:] //"   1920x1080   75   60" -> "1080   75"

	vres := line[:strings.Index(line, " ")] //"1080   75   60" -> "1080"

	refresh := line[strings.Index(line, " ")+1:] //"1080   75   60" -> "   75   60"

	if starIndex := strings.Index(refresh, "*"); starIndex > -1 {
		refresh = refresh[:starIndex]
		refresh = strings.ReplaceAll(refresh, " ", "")
	} else {
		for space := strings.Index(refresh, " "); space > -1; space = strings.Index(refresh, " ") {
			if space == 0 {
				refresh = refresh[1:]
			} else {
				refresh = refresh[:space]
			}
		} //"   75   60" -> "75"
	}

	refreshValue, err := strconv.ParseFloat(refresh, 64)

	if err != nil {
		return "", err
	}

	refreshValue = math.Ceil(refreshValue)
	return fmt.Sprintf("%vx%v@%v", hres, vres, refreshValue), nil
}

func (m *monitor) sortResolutions() []string {
	resolutions := make([]string, len(m.resolutions), len(m.resolutions))
	copy(resolutions, m.resolutions)
	pixels := make([]int64, len(resolutions))
	for i, v := range resolutions {
		hRes := v[:max(strings.Index(v, "x"), 0)]
		vRes := v[strings.Index(v, "x")+1 : strings.Index(v, "@")]
		hNum, err := strconv.ParseInt(hRes, 10, 32)
		vNum, err2 := strconv.ParseInt(vRes, 10, 32)
		if err != nil || err2 != nil {
			continue
		}
		pixels[i] = hNum * vNum
	}
	sortedNums := make([]int64, len(pixels), len(pixels))
	copy(sortedNums, pixels)
	slices.Sort(sortedNums)
	for i, sorted := range sortedNums {
		for j, unsorted := range pixels {
			if sorted == unsorted {
				tempN := pixels[i]
				pixels[i] = pixels[j]
				pixels[j] = tempN

				temp := resolutions[i]
				resolutions[i] = resolutions[j]
				resolutions[j] = temp
			}
		}
	}
	slices.Reverse(resolutions)
	m.resolutions = resolutions
	return resolutions
}
