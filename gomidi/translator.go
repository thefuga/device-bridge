package gomidi

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/thefuga/device-linker/footswitch"
	"github.com/thefuga/device-linker/keyboard"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewControlChangeSwitchTranslator,
	NewOutputDevice,
)

type (
	ControlChangeSwitchTranslator map[keyboard.Keypress]footswitch.Switch

	OutputDevice struct {
		port string
		out  drivers.Out
	}
)

func NewControlChangeSwitchTranslator() *ControlChangeSwitchTranslator {
	switches := viper.GetStringMap("gomidi.switches")
	translations := make(ControlChangeSwitchTranslator, len(switches))
	for k, v := range switches {
		s := v.(map[string]interface{})
		translations[keyboard.Keypress{Value: keyboard.InputValue(k)}] = footswitch.Switch{
			Channel:    uint8(s["channel"].(float64)),
			Controller: uint8(s["controller"].(float64)),
		}
	}
	return &translations
}

func NewOutputDevice() *OutputDevice {
	return &OutputDevice{
		port: viper.GetString("gomidi.port"),
	}
}

func (t ControlChangeSwitchTranslator) Translate(str keyboard.Keypress) (midi.Message, error) {
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
