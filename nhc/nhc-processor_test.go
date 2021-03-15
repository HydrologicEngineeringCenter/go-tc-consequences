package nhc

import (
	"testing"

	"github.com/USACE/go-consequences/hazards"
)

func Test_Convert(t *testing.T) {
	d := convertDepthtoHazardEvent(convertByteToDepth(1))
	de, _ := d.(hazards.DepthEvent)
	if de.Depth() != 1.0 {
		t.Errorf("expected 1.0 got something else.")
	}
	d2 := convertDepthtoHazardEvent(convertByteToDepth(2))
	de2, _ := d2.(hazards.DepthEvent)
	if de2.Depth() != 2.0 {
		t.Errorf("expected 2.0 got something else.")
	}
	d3 := convertDepthtoHazardEvent(convertByteToDepth(3))
	de3, _ := d3.(hazards.DepthEvent)
	if de3.Depth() != 3.0 {
		t.Errorf("expected 3.0 got something else.")
	}
	d4 := convertDepthtoHazardEvent(convertByteToDepth(4))
	de4, _ := d4.(hazards.DepthEvent)
	if de4.Depth() != 6.0 {
		t.Errorf("expected 6.0 got something else.")
	}
	d5 := convertDepthtoHazardEvent(convertByteToDepth(5))
	de5, _ := d5.(hazards.DepthEvent)
	if de5.Depth() != 9.0 {
		t.Errorf("expected 9.0 got something else.")
	}
}
