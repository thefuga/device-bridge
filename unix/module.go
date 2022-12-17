package unix

import (
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
	fx.Provide(
		newListener,
		newTranslator,
		newlinker,
		newLinkerFunc,
		newUnlinkerFunc,
		keyboard.NewStdIn,
		linker.NewSync,
	),
)

var Invokables = fx.Options(
	key2midi.Invokables,
)

type (
	UnixTranslator = linker.Translator[keyboard.Keypress, midi.Message]
	UnixLinker     = linker.Linker[keyboard.Keypress, midi.Message]
)

func newlinker(
	t *gomidi.ControlChangeSwitchTranslator,
	kb *keyboard.Keyboard,
	md *gomidi.OutputDevice,
	sync linker.Sync,
) *UnixLinker {
	return linker.NewLinker[keyboard.Keypress, midi.Message](t, kb, md, sync)
}

func newLinkerFunc(l *UnixLinker) linker.LinkFunc {
	return l.Link
}

func newUnlinkerFunc(l *UnixLinker) linker.UnlinkFunc {
	return l.Unlink
}

func newListener(stdin keyboard.StdIn) keyboard.Listener {
	return unix.NewListener(stdin)
}

func newTranslator(t gomidi.ControlChangeSwitchTranslator) UnixTranslator {
	return &t
}
