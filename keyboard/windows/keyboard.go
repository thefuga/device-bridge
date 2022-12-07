package windows

import (
	"context"
	"fmt"

	"github.com/thefuga/device-linker/keyboard"

	ghKeyboard "github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
)

var messageTypes = map[types.Message]int{
	types.WM_KEYDOWN: keyboard.KeyDown,
	types.WM_KEYUP:   keyboard.KeyUp,
}

type Listener struct{}

func NewListener() *Listener {
	return &Listener{}
}

func (*Listener) Listen(ctx context.Context, in chan keyboard.Keypress) {
	keyboardChan := make(chan types.KeyboardEvent, 100)
	ghKeyboard.Install(nil, keyboardChan)
	defer ghKeyboard.Uninstall()

	for {
		select {
		case <-ctx.Done():
			return
		case k := <-keyboardChan:
			in <- keyboard.Keypress{
				Value: keyboard.InputValue(fmt.Sprint(k.VKCode)),
				Type:  messageTypes[k.Message],
			}
		}
	}
}
