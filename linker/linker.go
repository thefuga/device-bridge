package linker

import (
	"context"
	"fmt"
)

type (
	Translator[In InputValue, Out OutputValue] interface {
		Translate(In) (Out, error)
	}

	InputValue interface {
		IsZero() bool
	}

	OutputValue interface{}

	InputDevice[In InputValue] interface {
		Listen(context.Context) error
		Process(context.Context) chan []In
	}

	OutputDevice[T OutputValue] interface {
		Send(T) error
	}

	Linker[In InputValue, Out OutputValue] struct {
		translator   Translator[In, Out]
		inputDevice  InputDevice[In]
		outputDevice OutputDevice[Out]
	}
)

func NewLinker[In InputValue, Out OutputValue](
	t Translator[In, Out], inDevice InputDevice[In], outDevice OutputDevice[Out],
) *Linker[In, Out] {
	return &Linker[In, Out]{
		translator:   t,
		inputDevice:  inDevice,
		outputDevice: outDevice,
	}
}

func (l *Linker[In, Out]) Link(parent context.Context) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	listenErr := l.inputDevice.Listen(ctx)

	if listenErr != nil {
		return listenErr
	}

	for {
		select {
		case <-parent.Done():
			return nil
		case inputs := <-l.inputDevice.Process(ctx):
			for _, input := range inputs {
				fmt.Printf("received input: %v\n", input)
				if err := l.translateAndSend(input); err != nil {
					fmt.Printf("translation error: %v", err) // TODO check error to see if linker must stop
				}
			}
		}
	}
}

func (b *Linker[In, Out]) translateAndSend(in In) error {
	message, translationErr := b.translator.Translate(in)

	if translationErr != nil {
		return translationErr
	}

	return b.outputDevice.Send(message)
}
