package unix

import (
	"context"
	"os/exec"

	"github.com/thefuga/device-linker/keyboard"
)

type (
	Listener struct {
		stdin keyboard.StdIn
	}
)

func NewListener(stdin keyboard.StdIn) *Listener {
	return &Listener{stdin: stdin}
}

func (l *Listener) Listen(ctx context.Context, in chan keyboard.Keypress) {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			b := make([]byte, 1)
			if n, err := l.stdin.Read(b); n < len(b) || err != nil {
				continue
			}
			in <- keyboard.Keypress{
				Value: keyboard.InputValue(b),
				Type:  keyboard.KeyDown,
			}
		}
	}
}
