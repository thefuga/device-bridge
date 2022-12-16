package footswitch

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thefuga/device-linker/footswitch"
	"github.com/thefuga/device-linker/keyboard"
	"github.com/thefuga/device-linker/linker"
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
	unlink     context.CancelFunc
	state      int
	stdin      keyboard.StdIn
	sync       linker.Sync
	once       sync.Once
}

func NewFootswitchScreen(
	f *footswitch.Footswitch,
	l linker.LinkFunc,
	stdin keyboard.StdIn,
	sync linker.Sync,
) FootswitchScreen {
	return FootswitchScreen{
		footswitch: f,
		link:       l,
		stdin:      stdin,
		sync:       sync,
	}
}

func (s FootswitchScreen) View() string {
	var renderedSwitches []string

	for _, v := range s.footswitch.Switches {
		renderedSwitches = append(
			renderedSwitches,
			fmt.Sprintf("%s  %s", v.Label, ledSymbol(v.On)+""),
		)
	}

	sort.Strings(renderedSwitches)

	return style.Render(strings.Join(renderedSwitches, "\n"))
}

func (m FootswitchScreen) Init() tea.Cmd {
	return nil
}

func (m FootswitchScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		m.unlink = cancel
		go m.link(ctx)
	})

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.unlink()
			m.once = sync.Once{}
		default:
			m.stdin.Write([]byte(msg.String()))
			select {
			case <-m.sync:
				return m, nil
			case <-time.After(1 * time.Second): // This prevents the screen from freezing in case the sync channel doesn't return
				return m, nil
			}
		}
	}

	return m, nil
}

func ledSymbol(s bool) string {
	if s {
		// return "🔴"
		return "◉"
	}
	return "○"
	// return "⚫"
}
