package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type resolutionSelectPage struct {
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

func (p resolutionSelectPage) New(name string, monitors map[string]*monitor) resolutionSelectPage {
	resolutions := monitors[name].resolutions
	return resolutionSelectPage{
		name:        name,
		monitor:     monitors[name],
		resolutions: resolutions,
		monitors:    monitors,
	}
}

func (m resolutionSelectPage) Init() tea.Cmd {
	return nil
}

func (m resolutionSelectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					len(m.monitor.modes[m.resolution])-1,
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
				m.refreshRate = m.monitor.modes[m.resolution][m.subcursor]
				err := setRes(
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
				return monitorSelectPage{}.New(m.monitors), nil
			}
			m.resolution = ""
		case "ctrl+c":
			return m, tea.Quit
		}

		m.previousInput = msg.String()
	}
	return m, nil
}

func (m resolutionSelectPage) View() string {
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

func (m resolutionSelectPage) renderRefreshRates() string {
	refreshRates := " Select A Refresh Rate\n" + "\n"
	for _, res := range m.resolutions {
		for i, rate := range m.monitor.modes[res] {

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
