package win

import (
	"context"

	"github.com/thefuga/device-linker/cmd/key2midi"
	"github.com/thefuga/device-linker/gomidi"
	"github.com/thefuga/device-linker/keyboard"
	"github.com/thefuga/device-linker/keyboard/windows"
	"github.com/thefuga/device-linker/linker"

	"gitlab.com/gomidi/midi/v2"
	"go.uber.org/fx"
)

var Module = fx.Options(
	key2midi.Module,
	fx.Provide(newListener, newTranslator, newlinker, linker.NewSync),
)

var Invokables = fx.Invoke(
	key2midi.Invokables,
	func(l *winLinker) error {
		return l.Link(context.Background())
	},
)

type (
	winTranslator = linker.Translator[keyboard.Keypress, midi.Message]
	winLinker     = linker.Linker[keyboard.Keypress, midi.Message]
)

func newlinker(
	t *gomidi.ControlChangeSwitchTranslator,
	kb *keyboard.Keyboard,
	md *gomidi.OutputDevice,
	sync linker.Sync,
) *winLinker {
	return linker.NewLinker[keyboard.Keypress, midi.Message](t, kb, md, sync)
}

func newListener() keyboard.Listener {
	return windows.NewListener()
}

func newTranslator(t gomidi.ControlChangeSwitchTranslator) winTranslator {
	return &t
}
