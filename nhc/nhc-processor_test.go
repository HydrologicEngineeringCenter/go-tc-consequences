package nhc

import (
	"testing"
)

func Test_Convert(t *testing.T) {
	d := convertDepthtoDepthEvent(convertByteToDepth(1))
	if d.Depth != 1.0 {
		t.Errorf("expected 1.0 got something else.")
	}
	d2 := convertDepthtoDepthEvent(convertByteToDepth(2))
	if d2.Depth != 2.0 {
		t.Errorf("expected 2.0 got something else.")
	}
	d3 := convertDepthtoDepthEvent(convertByteToDepth(3))
	if d3.Depth != 3.0 {
		t.Errorf("expected 3.0 got something else.")
	}
	d4 := convertDepthtoDepthEvent(convertByteToDepth(4))
	if d4.Depth != 6.0 {
		t.Errorf("expected 6.0 got something else.")
	}
	d5 := convertDepthtoDepthEvent(convertByteToDepth(5))
	if d5.Depth != 9.0 {
		t.Errorf("expected 9.0 got something else.")
	}
}
