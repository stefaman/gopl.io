package tempconv

import (
	"testing"
)

func TestTemp1(t *testing.T) {
	f := CToF(FreezingC)
	if f != 32 {
		t.Errorf("CToF failed. Got %v, but expert 32°F.", f)
	}

	c := FToC(FreezingF)
	if c != 0 {
		t.Errorf("FToC failed. Got %v, but expert 32°F.", c)
	}
}
