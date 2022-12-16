package home

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thefuga/device-linker/keyboard"
	"github.com/thefuga/device-linker/linker"
	"github.com/thefuga/device-linker/view/footswitch"
)

type HomeScreen struct {
	footswitch    footswitch.FootswitchScreen
	currentScreen tea.Model

	cursorPosition int
	options        []string

	stdin keyboard.StdIn
	sync  linker.Sync
}

var style = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(100)

var lineStyle = lipgloss.NewStyle().Align(lipgloss.Center)

func NewHomeScreen(f footswitch.FootswitchScreen) HomeScreen {
	m := HomeScreen{
		footswitch: f,
		options: []string{
			"  list keybinds",
			"  add new keybind",
			"  start listening",
		},
	}

	m.currentScreen = m

	return m
}

func (m HomeScreen) Init() tea.Cmd {
	return nil
}

func (m HomeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyMsg := msg.String()
		switch keyMsg {
		case "ctrl+c", "q":
			defer os.Exit(0)
			return m, tea.Quit
		case "esc":
			m.currentScreen, _ = m.currentScreen.Update(msg)
			m.currentScreen = m
		default:
			switch m.currentScreen.(type) {
			case HomeScreen:
				switch keyMsg {
				case "j":
					if m.cursorPosition >= len(m.options)-1 {
						m.cursorPosition = 0
					} else {
						m.cursorPosition += 1
					}
				case "k":
					if m.cursorPosition <= 0 {
						m.cursorPosition = len(m.options) - 1
					} else {
						m.cursorPosition -= 1
					}
				case "enter":
					switch m.cursorPosition {
					case 0:
						break
						// m.currentScreen = m.ListKeybindingsScreen
					case 1:
						break
						// m.currentScreen = m.AddNewKeybindScreen
					case 2:
						m.currentScreen = m.footswitch
					}
				}
			default:
				m.currentScreen, _ = m.currentScreen.Update(msg)
			}
		}
	}

	return m, nil
}

func (m HomeScreen) View() string {
	if _, ok := m.currentScreen.(HomeScreen); !ok {
		return m.currentScreen.View()
	}

	options := make([]string, len(m.options))
	copy(options, m.options)
	options[m.cursorPosition] = strings.Replace(options[m.cursorPosition], " ", ">", 1)

	return style.Render(strings.Join(options, "\n"))
}
