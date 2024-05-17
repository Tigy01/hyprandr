package main

import (
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type monitorSelectPage struct {
	selection    int
	monitors     map[string]*monitor
	monitorNames []string
}

func (p monitorSelectPage) New(monitors map[string]*monitor) monitorSelectPage {
	monitorNames := make([]string, 0)
	for n := range monitors {
		monitorNames = append(monitorNames, n)
	}
	slices.Sort(monitorNames)
	return monitorSelectPage{
		selection:    0,
		monitors:     monitors,
		monitorNames: monitorNames,
	}
}

func (m monitorSelectPage) Init() tea.Cmd {
	return nil
}

func (m monitorSelectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		windowWidth = msg.Width
		windowHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			m.selection = max(m.selection-1, 0)
			return m, nil
		case "j":
			m.selection = min(m.selection+1, len(m.monitorNames)-1)
			return m, nil
		case "enter":
			nextPage := resolutionSelectPage{}.New(
				m.monitorNames[m.selection],
				m.monitors,
			)
			return nextPage, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m monitorSelectPage) View() string {
	output := ""
	for i, name := range m.monitorNames {
		output += "\n"
		if i == m.selection {
			output += ">[" + name + "]"
		} else {
			output += "  " + name
		}
	}
	return lipgloss.Place(
		windowWidth,
		windowHeight,
		lipgloss.Left,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			output,
			lipgloss.PlaceVertical(
				windowHeight-(len(m.monitorNames)),
				lipgloss.Bottom,
				lipgloss.NewStyle().Render(getHelp()),
			),
		),
	)
}
