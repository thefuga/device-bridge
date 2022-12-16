package view

import (
	"github.com/thefuga/device-linker/view/footswitch"
	"github.com/thefuga/device-linker/view/home"
	"go.uber.org/fx"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	Module = fx.Provide(
		footswitch.NewFootswitchScreen,
		home.NewHomeScreen,
		NewProgram,
	)

	Invokables = fx.Invoke(RunProgram)
)

func NewProgram(home home.HomeScreen) *tea.Program {
	return tea.NewProgram(home)
}

func RunProgram(p *tea.Program) error {
	_, err := p.Run()
	return err
}
