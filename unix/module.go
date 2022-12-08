package unix

import (
	"context"

	"github.com/thefuga/device-linker/cmd/key2midi"
	"github.com/thefuga/device-linker/gomidi"
	"github.com/thefuga/device-linker/keyboard"
	"github.com/thefuga/device-linker/keyboard/unix"
	"github.com/thefuga/device-linker/linker"

	"gitlab.com/gomidi/midi/v2"
	"go.uber.org/fx"
)

var Module = fx.Options(
	key2midi.Module,
	fx.Provide(newListener, newTranslator, newlinker),
)

var Invokables = fx.Invoke(
	func(l *UnixLinker) error {
		return l.Link(context.Background())
	},
)

type (
	UnixTranslator = linker.Translator[keyboard.Keypress, midi.Message]
	UnixLinker     = linker.Linker[keyboard.Keypress, midi.Message]
)

func newlinker(
	t *gomidi.ControlChangeSwitchTranslator,
	kb *keyboard.Keyboard,
	md *gomidi.OutputDevice,
) *UnixLinker {
	return linker.NewLinker[keyboard.Keypress, midi.Message](t, kb, md)
}

func newListener() keyboard.Listener {
	return unix.NewListener()
}

func newTranslator(t gomidi.ControlChangeSwitchTranslator) UnixTranslator {
	return &t
}
