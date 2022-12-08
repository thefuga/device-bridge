package gomidi

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/thefuga/device-linker/keyboard"
	"gitlab.com/gomidi/midi/v2"
)

func TestTranslate(t *testing.T) {
	viper.Set("gomidi.switches", map[string]interface{}{
		"a": map[string]interface{}{
			"channel":    1.,
			"controller": 1.,
		},
	})
	defer viper.Reset()

	expectedFirstPress := midi.Message([]byte{0b10110001, 0b1, 0b1111111})
	expectedSecondPress := midi.Message([]byte{0b10110001, 0b1, 0b000000})

	translator := NewControlChangeSwitchTranslator()

	keypress := keyboard.Keypress{Value: keyboard.InputValue("a"), Type: keyboard.KeyDown}

	midiMessage, _ := translator.Translate(keypress)
	if !reflect.DeepEqual(midiMessage, expectedFirstPress) {
		t.Errorf("expected %s, got %s", expectedFirstPress, midiMessage)
	}

	midiMessage, _ = translator.Translate(keypress)
	if !reflect.DeepEqual(midiMessage, expectedSecondPress) {
		t.Errorf("expected %s, got %s", expectedSecondPress.String(), midiMessage.String())
	}
}

func TestTranslate_InvalidKey(t *testing.T) {
	translator := NewControlChangeSwitchTranslator()
	keypress := keyboard.Keypress{Value: keyboard.InputValue("a"), Type: keyboard.KeyDown}

	if _, err := translator.Translate(keypress); err == nil {
		t.Error("expected error from invalid translation")
	}
}
