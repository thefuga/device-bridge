package footswitch

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thefuga/device-linker/footswitch"
	"github.com/thefuga/device-linker/keyboard"
	"github.com/thefuga/device-linker/linker"
	"github.com/thefuga/go-collections"
)

var style = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(100)

type FootswitchScreen struct {
	footswitch *footswitch.Footswitch
	link       linker.LinkFunc
	unlink     linker.UnlinkFunc
	state      int
	stdin      keyboard.StdIn
	sync       linker.Sync
}

func NewFootswitchScreen(
	f *footswitch.Footswitch,
	l linker.LinkFunc,
	ul linker.UnlinkFunc,
	stdin keyboard.StdIn,
	sync linker.Sync,
) FootswitchScreen {
	return FootswitchScreen{
		footswitch: f,
		link:       l,
		unlink:     ul,
		stdin:      stdin,
		sync:       sync,
	}
}

func (s FootswitchScreen) View() string {
	renderedSwitches := collections.Map(s.footswitch.Switches, serializeSwitch)
	sort.Strings(renderedSwitches)

	return style.Render(strings.Join(renderedSwitches, "\n"))
}

func serializeSwitch(_ int, s *footswitch.Switch) string {
	return fmt.Sprintf("%s  %s", s.Label, ledSymbol(s.On)+"")
}

func (m FootswitchScreen) Init() tea.Cmd {
	return nil
}

func (m FootswitchScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.link(context.Background())

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.updateKeyMessage(msg)
	}

	return m, nil
}

func (m FootswitchScreen) updateKeyMessage(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.unlink()
	default:
		m.processInput(msg.String())
	}

	return m, nil
}

func (m *FootswitchScreen) processInput(msg string) {
	m.stdin.Write([]byte(msg))
	select {
	case <-m.sync:
		return
	case <-time.After(1 * time.Second): // This prevents the screen from freezing in case the sync channel doesn't return
		return
	}
}

func ledSymbol(s bool) string {
	if s {
		// return "🔴"
		return "◉"
	}
	return "○"
	// return "⚫"
}
