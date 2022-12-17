package linker

import (
	"context"
	"sync"
)

type (
	LinkFunc   func(context.Context) error
	UnlinkFunc func()

	Sync chan struct{}

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
		sync         Sync
		once         sync.Once
		unlink       context.CancelFunc
	}
)

func NewLinker[In InputValue, Out OutputValue](
	t Translator[In, Out], inDevice InputDevice[In], outDevice OutputDevice[Out], sync Sync,
) *Linker[In, Out] {
	return &Linker[In, Out]{
		translator:   t,
		inputDevice:  inDevice,
		outputDevice: outDevice,
		sync:         sync,
	}
}

func NewSync() Sync {
	return make(Sync)
}

func (l *Linker[In, Out]) Unlink() {
	if l.unlink != nil {
		l.unlink()
	}

	l.once = sync.Once{}
}

func (l *Linker[In, Out]) Link(parent context.Context) (err error) {
	l.once.Do(func() {
		ctx, cancel := context.WithCancel(parent)
		l.unlink = cancel

		if err = l.listenInput(ctx); err != nil {
			return
		}

		go l.link(ctx)

	})

	return
}

func (l *Linker[In, Out]) listenInput(ctx context.Context) error {
	// ctx, cancel := context.WithCancel(parent)
	// defer cancel()

	listenErr := l.inputDevice.Listen(ctx)

	return listenErr
}

func (l *Linker[In, Out]) link(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	for {
		select {
		case <-parent.Done():
			return
		case inputs := <-l.inputDevice.Process(ctx):
			for _, input := range inputs {
				if err := l.translateAndSend(input); err != nil {
					// TODO  handle this error
				}
			}
		}
	}
}

func (b *Linker[In, Out]) translateAndSend(in In) error {
	defer b.syncState()

	message, translationErr := b.translator.Translate(in)

	if translationErr != nil {
		return translationErr
	}

	return b.outputDevice.Send(message)
}

func (b *Linker[In, Out]) syncState() {
	b.sync <- struct{}{}
}
