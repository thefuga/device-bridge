package key2midi

import (
	"github.com/thefuga/device-linker/footswitch"
	"github.com/thefuga/device-linker/gomidi"
	"github.com/thefuga/device-linker/keyboard"
	"github.com/thefuga/device-linker/view"

	"go.uber.org/fx"
)

var Module = fx.Options(
	keyboard.Module,
	footswitch.Module,
	gomidi.Module,
	view.Module,
)

var Invokables = fx.Options(view.Invokables)
