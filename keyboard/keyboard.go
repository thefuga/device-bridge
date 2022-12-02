package keyboard

import (
	"context"
	"os"
	"os/exec"
	"time"

	"go.uber.org/fx"
)

var Module = fx.Provide(NewKeyboard)

type (
	Keyboard struct {
		in            chan string
		out           chan string
		DebounceDelay time.Duration
	}

	InputValue string
)

func NewKeyboard(delay time.Duration) *Keyboard {
	return &Keyboard{
		DebounceDelay: delay,
	}
}

func (v InputValue) IsZero() bool {
	return v == ""
}

func (kb *Keyboard) Listen(ctx context.Context) {
	if kb.in == nil {
		kb.in = make(chan string)
	}

	go kb.listen(ctx)
}

func (kb *Keyboard) listen(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		// disable input buffering
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		// do not display entered characters on the screen
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			kb.in <- string(b)
		}
	}
}

func (kb *Keyboard) Process(ctx context.Context) chan string {
	if kb.out == nil {
		kb.out = make(chan string)
	}

	go kb.processKeystrokes(ctx)
	return kb.out
}

func (kb *Keyboard) processKeystrokes(ctx context.Context) {
	var buffer string
	for {
		select {
		case stdin, _ := <-kb.in:
			if stdin != "" {
				buffer += stdin
				continue
			}
		case <-time.After(kb.DebounceDelay):
			break
		case <-ctx.Done():
			return
		}

		break
	}

	kb.out <- buffer
}
