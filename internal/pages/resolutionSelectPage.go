package pages 

import (
	"fmt"

    "github.com/Tigy01/hyprandr/internal/cli"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResolutionSelectPage struct {
	cursor        int
	subcursor     int
	previousInput string
	resolution    string
	refreshRate   int

	name        string
	resolutions []string

	monitor  *monitor
	monitors map[string]*monitor
}

func (p ResolutionSelectPage) New(name string, monitors map[string]*monitor) ResolutionSelectPage {
	resolutions := monitors[name].Resolutions
	return ResolutionSelectPage{
		name:        name,
		monitor:     monitors[name],
		resolutions: resolutions,
		monitors:    monitors,
	}
}

func (m ResolutionSelectPage) Init() tea.Cmd {
	return nil
}

func (m ResolutionSelectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		windowWidth = msg.Width
		windowHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			if m.resolution != "" {
				m.subcursor = max(m.subcursor-1, 0)
			}
		case "j":
			if m.resolution == "" {
				m.cursor = min(m.cursor+1, len(m.resolutions)-1)
			}
		case "k":
			if m.resolution == "" {
				m.cursor = max(m.cursor-1, 0)
			}
		case "l":
			if m.resolution != "" {
				m.subcursor = min(
					m.subcursor+1,
					len(m.monitor.Modes[m.resolution])-1,
				)
			}
		case "g":
			if m.resolution != "" {
				break
			}
			if m.previousInput == "g" {
				m.cursor = 0
			}
		case "G":
			if m.resolution == "" {
				m.cursor = len(m.resolutions) - 1
			}
		case "enter":
			if m.resolution == "" {
				m.resolution = m.resolutions[m.cursor]
				m.subcursor = 0
			} else {
				m.refreshRate = m.monitor.Modes[m.resolution][m.subcursor]
				err := cli.SetRes(
					m.monitors,
					m.name,
					fmt.Sprintf("%v@%v", m.resolution, m.refreshRate),
				)

				if err != nil {
					panic(fmt.Sprintf("err: %v\n", err))
				}

				return m, nil
			}
		case "q":
			if m.resolution == "" {
				return MonitorSelectPage{}.New(m.monitors), nil
			}
			m.resolution = ""
		case "ctrl+c":
			return m, tea.Quit
		}

		m.previousInput = msg.String()
	}
	return m, nil
}

func (m ResolutionSelectPage) View() string {
	output := fmt.Sprintf(" [[%v]]\n\n", m.name)
	for i, res := range m.resolutions {

		if i == m.cursor {
			output += ">["
			output += fmt.Sprint(res, "]\n")
		} else {
			output += "  "
			output += fmt.Sprint(res, " \n")
		}

	}

	var refreshRates string
	if m.resolution != "" {
		refreshRates = m.renderRefreshRates()
	}

	withRightBorder := lipgloss.NewStyle().
    PaddingRight(3).
		BorderStyle(lipgloss.NormalBorder()).
		BorderRight(true)

	return lipgloss.Place(
		windowWidth,
		windowHeight,
		lipgloss.Left,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,

			lipgloss.JoinHorizontal(
				lipgloss.Top,
				withRightBorder.Render(output),
				refreshRates,
			),
			lipgloss.PlaceVertical(
				windowHeight-(len(m.resolutions)+3),
				lipgloss.Bottom,
				getHelp(),
			),
		),
	)
}

func (m ResolutionSelectPage) renderRefreshRates() string {
	refreshRates := " Select A Refresh Rate\n" + "\n"
	for _, res := range m.resolutions {
		for i, rate := range m.monitor.Modes[res] {

			if res != m.resolution {
				continue
			}

			if i == m.subcursor {
				refreshRates += fmt.Sprint(" >[", rate, "] ")
			} else {
				refreshRates += fmt.Sprint("   ", rate, "  ")
			}
		}

		refreshRates += "\n"
	}
	return refreshRates
}
