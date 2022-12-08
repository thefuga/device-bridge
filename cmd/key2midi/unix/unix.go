package main

import (
	"github.com/thefuga/device-linker/cmd/key2midi/config"
	"github.com/thefuga/device-linker/unix"

	"go.uber.org/fx"
)

func main() {
	config.Load()

	fx.New(unix.Module, unix.Invokables).Run()
}
