package footswitch

import "go.uber.org/fx"

var Module = fx.Provide(NewFootswitch)

type (
	Footswitch struct {
		Switches []Switch `json:"switches"`
	}

	Switch struct {
		Channel    uint8 `json:"channel"`
		Controller uint8 `json:"controller"`
		On         bool  `json:"-"`
	}
)

func NewFootswitch() *Footswitch {
	return &Footswitch{}
}

func (k *Switch) Press() *Switch {
	k.On = !k.On
	return k
}

func (k *Switch) Value() uint8 {
	if k.On {
		return 255
	}

	return 0
}
