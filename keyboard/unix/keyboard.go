package unix

import (
	"context"
	"os"
	"os/exec"

	"github.com/thefuga/device-linker/keyboard"
)

type Listener struct{}

func NewListener() *Listener {
	return &Listener{}
}

func (*Listener) Listen(ctx context.Context, in chan keyboard.Keypress) {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	var b []byte = make([]byte, 1)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			os.Stdin.Read(b)
			in <- keyboard.Keypress{
				Value: keyboard.InputValue(b),
				Type:  keyboard.KeyDown,
			}
		}
	}
}
