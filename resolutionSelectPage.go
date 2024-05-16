package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type resolutionSelectPage struct {
	name        string
	monitors    map[string]*monitor
	resolutions []string
	selection   int
}

func (p resolutionSelectPage) New(name string, monitors map[string]*monitor) resolutionSelectPage {
	resolutions := monitors[name].resolutions
	return resolutionSelectPage{
		name:        name,
		resolutions: resolutions,
		monitors:    monitors,
	}
}

func (m resolutionSelectPage) Init() tea.Cmd {
	return nil
}

func (m resolutionSelectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			m.selection = max(m.selection-1, 0)
		case "j":
			m.selection = min(m.selection+1, len(m.resolutions))
		case "q":
			return monitorSelectPage{}.New(m.monitors), nil
		case "enter":
			if m.selection == len(m.resolutions) {
				m.monitors = getResolutions(m.monitors)
				rewriteConfig(m.monitors)
                return m.New(m.name,m.monitors), nil
			}
			changeRes(m.monitors, m.name, m.selection)
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m resolutionSelectPage) View() string {
	output := fmt.Sprintf(" [[%v]]\n\n", m.name)
	for i, res := range m.resolutions {
		if i == m.selection {
			output += ">["
			output += fmt.Sprint(res, "]\n")
		} else {
			output += "  "
			output += fmt.Sprint(res, "\n")
		}
	}
	if m.selection == len(m.resolutions) {
		output += "\n>[[REFRESH RESOLUTIONS]]"
	} else {
		output += "\n  [REFRESH RESOLUTIONS] "
	}
    output += "-WARNING: This can cause resolutions to disappear"
	return output
}
