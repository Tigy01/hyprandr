package pages

import (
	"fmt"
	"slices"

	"github.com/Tigy01/hyprandr/internal/cli"
	"github.com/Tigy01/hyprandr/internal/monitors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type monitor = monitors.Monitor

type MonitorSelectPage struct {
	cursor        int
	monitors      map[string]*monitor
	monitorNames  []string
	previousInput string
	numDisabled   int
}

var windowWidth int
var windowHeight int

func getHelp() string {
	return lipgloss.JoinHorizontal(lipgloss.Bottom,
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true, false, true, true).
			Foreground(lipgloss.Color("#5555ff")).
			Render(" HELP "+lipgloss.NormalBorder().Right),
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true, true, true, false).
			Render(" h -> Left | j -> Down | k -> Up | l -> Right | RETURN -> Select | q -> Back | r -> Refresh | d -> Disable"),
	)
}

func (page MonitorSelectPage) New(monitors map[string]*monitor) MonitorSelectPage {
	numDisabled := 0
	monitorNames := make([]string, 0)
	for name, monitor := range monitors {
		if monitor.Disable {
			numDisabled += 1
		}
		monitorNames = append(monitorNames, name)
	}
	slices.Sort(monitorNames)
	return MonitorSelectPage{
		cursor:       0,
		monitors:     monitors,
		monitorNames: monitorNames,
		numDisabled:  numDisabled,
	}
}

func (page MonitorSelectPage) Init() tea.Cmd {
	return nil
}

func (page MonitorSelectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		windowWidth = msg.Width
		windowHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			page.cursor = max(page.cursor-1, 0)
		case "j":
			page.cursor = min(page.cursor+1, len(page.monitorNames)-1)
		case "g":
			if page.previousInput == "g" {
				page.cursor = 0
			}
		case "G":
			page.cursor = len(page.monitorNames) - 1
		case "enter":
			nextPage := ResolutionSelectPage{}.New(
				page.monitorNames[page.cursor],
				page.monitors,
			)
			return nextPage, nil
		case "ctrl+c", "q":
			return page, tea.Quit
		case "r":
			var err error
			page.monitors, err = cli.GetCurrentSettings()
			if err != nil {
				fmt.Println(err)
				return nil, tea.Quit
			}
			cli.RewriteConfig(page.monitors)
			return MonitorSelectPage{}.New(page.monitors), nil
		case "d":
			currentMonitor := page.monitors[page.monitorNames[page.cursor]]
			currentMonitor.Disable = !currentMonitor.Disable
			cli.RewriteConfig(page.monitors)
			return MonitorSelectPage{}.New(page.monitors), nil
		}
		page.previousInput = msg.String()
	}
	return page, nil
}

func padRight(input string, lineLength int, fill string) string {
	for len(input) < lineLength {
		input += fill
	}
	return input
}

func truncate(input string, lineLength int) string {
	if len(input) > lineLength {
		return input[0:lineLength-2] + "..."
	}
	return input
}

func (page MonitorSelectPage) View() string {
	var names string
	var resolutions string
	for i, name := range page.monitorNames {
		var line string
		if i == page.cursor {
			line += ">[" + name + "]"
		} else {
			line += "  " + name + " "
		}
		line = padRight(line, 12, " ")
		line = truncate(line, 12)
		names += line + "\n"
		resolutions += page.monitors[name].CurrentRes
		if page.monitors[name].Disable {
			resolutions += " #Disabled"
		}
		resolutions += "\n"
	}
	names = lipgloss.NewStyle().
		Border(
			lipgloss.NormalBorder(),
			false,
			true,
			false,
			false,
		).MarginRight(1).
		Render(names)
	return lipgloss.Place(
		windowWidth,
		windowHeight,
		lipgloss.Left,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				names,
				resolutions,
			),
			lipgloss.PlaceVertical(
				windowHeight-(len(page.monitorNames)+1),
				lipgloss.Bottom,
				lipgloss.NewStyle().Render(getHelp()),
			),
		),
	)
}
