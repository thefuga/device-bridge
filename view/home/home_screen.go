package home

import (
	"os"
	"strings"

	"github.com/thefuga/device-linker/view/footswitch"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeScreen struct {
	screens        []tea.Model
	currentScreen  tea.Model
	cursorPosition int
	options        []string
}

var style = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(100)

func NewHomeScreen(f footswitch.FootswitchScreen) HomeScreen {
	m := HomeScreen{
		screens: []tea.Model{nil, nil, f},
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
		return m.updateKeyMessage(msg)
	}

	return m, nil
}

func (m HomeScreen) View() string {
	if m.isAtHomeScreen() {
		return m.viewHomeScreen()
	}

	return m.currentScreen.View()
}

func (m HomeScreen) isAtHomeScreen() bool {
	_, ok := m.currentScreen.(HomeScreen)
	return ok
}

func (m HomeScreen) viewHomeScreen() string {
	options := make([]string, len(m.options))
	copy(options, m.options)
	options[m.cursorPosition] = strings.Replace(options[m.cursorPosition], " ", ">", 1)

	return style.Render(strings.Join(options, "\n"))
}

func (m *HomeScreen) updateKeyMessage(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.exit()
	case "esc":
		m.updateCurrentScreenAndReturnHome(msg)
	default:
		m.redirectMessage(msg)
	}

	return *m, nil
}

func (m *HomeScreen) exit() {
	os.Exit(0)
}

func (m *HomeScreen) updateCurrentScreenAndReturnHome(msg tea.KeyMsg) {
	m.updateCurrentScreen(msg)
	m.currentScreen = *m
}

func (m *HomeScreen) updateCurrentScreen(msg tea.Msg) {
	m.currentScreen, _ = m.currentScreen.Update(msg)
}

func (m *HomeScreen) redirectMessage(msg tea.KeyMsg) {
	switch m.currentScreen.(type) {
	case HomeScreen:
		m.updateMenu(msg)
	default:
		m.updateCurrentScreen(msg)
	}
}

func (m *HomeScreen) updateMenu(msg tea.KeyMsg) {
	switch msg.String() {
	case "j":
		m.moveCursorDown()
	case "k":
		m.moveCursorUp()
	case "enter":
		m.selectScreen()
	}
}

func (m *HomeScreen) selectScreen() {
	m.currentScreen = m.screens[m.cursorPosition]
}

func (m *HomeScreen) moveCursorDown() {
	if m.cursorPosition >= len(m.options)-1 {
		m.cursorPosition = 0
	} else {
		m.cursorPosition += 1
	}
}

func (m *HomeScreen) moveCursorUp() {
	if m.cursorPosition <= 0 {
		m.cursorPosition = len(m.options) - 1
	} else {
		m.cursorPosition -= 1
	}
}
