package keyboard

import (
	"bytes"
	"context"
	"io"
	"time"

	"go.uber.org/fx"

	"github.com/spf13/viper"
)

var Module = fx.Provide(NewKeyboard)

const (
	KeyDown = iota
	KeyUp
)

type (
	Listener interface {
		Listen(context.Context, chan Keypress)
	}

	Keyboard struct {
		in       chan Keypress
		out      chan []Keypress
		listener Listener

		DebounceDelay time.Duration
	}

	Keypress struct {
		Value InputValue
		Type  int
	}

	InputValue string

	StdIn interface {
		io.Reader
		io.Writer
		Reset()
	}
)

func NewStdIn() StdIn {
	return &bytes.Buffer{}

}

func NewKeyboard(listener Listener) *Keyboard {
	return &Keyboard{
		DebounceDelay: viper.GetDuration("keyboard.debaunce_delay"),
		listener:      listener,
	}
}

func (k Keypress) IsZero() bool {
	return k.Value == ""
}

func (kb *Keyboard) Listen(ctx context.Context) error {
	if kb.in == nil {
		kb.in = make(chan Keypress)
	}

	go kb.listener.Listen(ctx, kb.in)
	return nil
}

func (kb *Keyboard) Process(ctx context.Context) chan []Keypress {
	if kb.out == nil {
		kb.out = make(chan []Keypress)
	}

	go kb.processKeystrokes(ctx)
	return kb.out
}

func (kb *Keyboard) processKeystrokes(ctx context.Context) {
	var buffer []Keypress

	select {
	case stdin, _ := <-kb.in:
		if stdin.Value != "" {
			buffer = append(buffer, stdin)
		}
	case <-ctx.Done():
		return
	}

	if buffer == nil {
		kb.out <- buffer
		return
	}

	kb.out <- buffer[:1]
}
