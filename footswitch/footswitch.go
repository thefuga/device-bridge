package footswitch

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewFootswitch)

type (
	Footswitch struct {
		Switches []*Switch `json:"switches"`
	}

	Switch struct {
		Channel    uint8 `json:"channel"`
		Controller uint8 `json:"controller"`
		On         bool  `json:"-"`
		Label      string
	}
)

func NewFootswitch() *Footswitch {
	configSwitches := viper.GetStringMap("gomidi.switches")
	var switches []*Switch

	for k, v := range configSwitches {
		s := v.(map[string]interface{})

		switches = append(switches, &Switch{
			Channel:    uint8(s["channel"].(float64)),
			Controller: uint8(s["controller"].(float64)),
			Label:      k,
		})
	}

	return &Footswitch{Switches: switches}
}

func (k *Switch) Press() *Switch {
	k.On = !k.On
	return k
}

func (k *Switch) IsOn() bool {
	return k.On
}

func (k *Switch) Value() uint8 {
	if k.On {
		return 255
	}

	return 0
}
