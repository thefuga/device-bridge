package footswitch

import (
	"github.com/spf13/viper"
	"github.com/thefuga/go-collections"
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
	return &Footswitch{Switches: buildGomidiSwitches()}
}

func buildGomidiSwitches() []*Switch {
	return collections.Map(
		getSlice("gomidi.switches"),
		interfaceToSwitch,
	)
}

func getSlice(key string) []interface{} {
	return viper.Get(key).([]interface{})
}

func interfaceToSwitch(_ int, v interface{}) *Switch {
	s := v.(map[string]interface{})
	return &Switch{
		Channel:    uint8(s["channel"].(float64)),
		Controller: uint8(s["controller"].(float64)),
		Label:      s["label"].(string),
	}
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
