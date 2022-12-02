//go:build windows && gomidi
// +build windows,gomidi

package main

import (
	"github.com/thefuga/device-linker/footswitch"
	"github.com/thefuga/device-linker/gomidi"
	"github.com/thefuga/device-linker/keyboard"
	"github.com/thefuga/device-linker/linker"

	"gitlab.com/gomidi/midi/v2"
	"go.uber.org/fx"
)

var Module = fx.Options(
	keyboard.Module,
	footswitch.Module,
	gomidi.Module,
	fx.Provide(linker.NewLinker[keyboard.InputValue, midi.Message]),
)

func main() {
	fx.New(Module).Run()
}
