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

func (page ResolutionSelectPage) New(name string, monitors map[string]*monitor) ResolutionSelectPage {
	resolutions := monitors[name].Resolutions
	return ResolutionSelectPage{
		name:        name,
		monitor:     monitors[name],
		resolutions: resolutions,
		monitors:    monitors,
	}
}

func (page ResolutionSelectPage) Init() tea.Cmd {
	return nil
}

func (page ResolutionSelectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		windowWidth = msg.Width
		windowHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			if page.resolution != "" {
				page.subcursor = max(page.subcursor-1, 0)
			}
		case "j":
			if page.resolution == "" {
				page.cursor = min(page.cursor+1, len(page.resolutions)-1)
			}
		case "k":
			if page.resolution == "" {
				page.cursor = max(page.cursor-1, 0)
			}
		case "l":
			if page.resolution != "" {
				page.subcursor = min(
					page.subcursor+1,
					len(page.monitor.Modes[page.resolution])-1,
				)
			}
		case "g":
			if page.resolution != "" {
				break
			}
			if page.previousInput == "g" {
				page.cursor = 0
			}
		case "G":
			if page.resolution == "" {
				page.cursor = len(page.resolutions) - 1
			}
		case "enter":
			if page.resolution == "" {
				page.resolution = page.resolutions[page.cursor]
				page.subcursor = 0
			} else {
				page.refreshRate = page.monitor.Modes[page.resolution][page.subcursor]
				err := cli.SetRes(
					page.monitors,
					page.name,
					fmt.Sprintf("%v@%v", page.resolution, page.refreshRate),
				)

				if err != nil {
					fmt.Println(err)
					return nil, tea.Quit
				}

				return page, nil
			}
		case "q":
			if page.resolution == "" {
				return MonitorSelectPage{}.New(page.monitors), nil
			}
			page.resolution = ""
		case "ctrl+c":
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
		}

		page.previousInput = msg.String()
	}
	return page, nil
}

func (page ResolutionSelectPage) View() string {
	output := fmt.Sprintf(" [[%v]]\n\n", page.name)
	for i, res := range page.resolutions {

		if i == page.cursor {
			output += ">["
			output += fmt.Sprint(res, "]\n")
		} else {
			output += "  "
			output += fmt.Sprint(res, " \n")
		}

	}

	var refreshRates string
	if page.resolution != "" {
		refreshRates = page.renderRefreshRates()
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
				windowHeight-(len(page.resolutions)+3),
				lipgloss.Bottom,
				getHelp(),
			),
		),
	)
}

func (page ResolutionSelectPage) renderRefreshRates() string {
	refreshRates := " Select A Refresh Rate\n" + "\n"
	for _, res := range page.resolutions {
		for i, rate := range page.monitor.Modes[res] {

			if res != page.resolution {
				continue
			}

			if i == page.subcursor {
				refreshRates += fmt.Sprint(" >[", rate, "] ")
			} else {
				refreshRates += fmt.Sprint("   ", rate, "  ")
			}
		}

		refreshRates += "\n"
	}
	return refreshRates
}
