package main

import (
	"github.com/thefuga/device-linker/cmd/key2midi/config"
	"github.com/thefuga/device-linker/win"

	"go.uber.org/fx"
)

func main() {
	config.Load()

	fx.New(win.Module, win.Invokables).Run()
}
