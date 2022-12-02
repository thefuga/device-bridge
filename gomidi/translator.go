package gomidi

import (
	"fmt"

	"github.com/thefuga/device-linker/footswitch"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewControlChangeSwitchTranslator,
	NewOutputDevice,
)

type (
	ControlChangeSwitchTranslator map[string]footswitch.Switch

	OutputDevice struct {
		port string
		out  drivers.Out
	}
)

func NewControlChangeSwitchTranslator() *ControlChangeSwitchTranslator {
	return &ControlChangeSwitchTranslator{}
}

func NewOutputDevice(port string) *OutputDevice {
	return &OutputDevice{
		port: port,
	}
}

func (t ControlChangeSwitchTranslator) Translate(str string) (midi.Message, error) {
	s, ok := t[str]

	if !ok {
		return nil, fmt.Errorf("not found")
	}

	return midi.ControlChange(s.Channel, s.Controller, s.Press().Value()), nil
}

func (out OutputDevice) Send(message midi.Message) error {
	send, err := midi.SendTo(out.out)
	if err != nil {
		return err
	}

	return send(message)
}
