package gomidi

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/thefuga/device-linker/footswitch"
	"github.com/thefuga/device-linker/keyboard"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/midicatdrv"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewControlChangeSwitchTranslator,
	NewOutputDevice,
)

type (
	ControlChangeSwitchTranslator map[keyboard.Keypress]*footswitch.Switch

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
		translations[keyboard.Keypress{Value: keyboard.InputValue(k)}] = &footswitch.Switch{
			Channel:    uint8(s["channel"].(float64)),
			Controller: uint8(s["controller"].(float64)),
		}
	}

	return &translations
}

func NewOutputDevice() *OutputDevice {
	out, err := drivers.OutByName(viper.GetString("gomidi.out_port_name"))
	if err != nil {
		panic(err)
	}
	return &OutputDevice{
		out: out,
	}
}

func (t *ControlChangeSwitchTranslator) Translate(str keyboard.Keypress) (midi.Message, error) {
	if str.Type == keyboard.KeyUp {
		return nil, fmt.Errorf("key up action not defined")
	}

	s, ok := (*t)[keyboard.Keypress{Value: str.Value}]

	if !ok {
		return nil, fmt.Errorf("translation not found")
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
