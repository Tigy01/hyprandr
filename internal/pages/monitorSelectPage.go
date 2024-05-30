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
			Render(" h -> Left | j -> Down | k -> Up | l -> Right | RETURN -> Select | q -> Back | r -> Refresh "),
	)
}

func (p MonitorSelectPage) New(monitors map[string]*monitor) MonitorSelectPage {
	monitorNames := make([]string, 0)
	for n := range monitors {
		monitorNames = append(monitorNames, n)
	}
	slices.Sort(monitorNames)
	return MonitorSelectPage{
		cursor:       0,
		monitors:     monitors,
		monitorNames: monitorNames,
	}
}

func (m MonitorSelectPage) Init() tea.Cmd {
	return nil
}

func (m MonitorSelectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		windowWidth = msg.Width
		windowHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			m.cursor = max(m.cursor-1, 0)
		case "j":
			m.cursor = min(m.cursor+1, len(m.monitorNames)-1)
		case "g":
			if m.previousInput == "g" {
				m.cursor = 0
			}
		case "G":
			m.cursor = len(m.monitorNames) - 1
		case "enter":
			nextPage := ResolutionSelectPage{}.New(
				m.monitorNames[m.cursor],
				m.monitors,
			)
			return nextPage, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			var err error
			m.monitors, err = cli.GetCurrentSettings()
			if err != nil {
				fmt.Println(err)
				return nil, tea.Quit
			}
            cli.RewriteConfig(m.monitors)
			return MonitorSelectPage{}.New(m.monitors), nil
		}
		m.previousInput = msg.String()
	}
	return m, nil
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

func (m MonitorSelectPage) View() string {
	var names string
	var resolutions string
	for i, name := range m.monitorNames {
		var line string
		if i == m.cursor {
			line += ">[" + name + "]"
		} else {
			line += "  " + name + " "
		}
		line = padRight(line, 12, " ")
		line = truncate(line, 12)
		names += line + "\n"
		resolutions += m.monitors[name].CurrentRes + "\n"
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
				windowHeight-(len(m.monitorNames)+1),
				lipgloss.Bottom,
				lipgloss.NewStyle().Render(getHelp()),
			),
		),
	)
}
