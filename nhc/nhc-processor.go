package nhc

import (
	"github.com/USACE/go-consequences/hazards"
)

func convertDepthtoDepthEvent(d float64) hazards.DepthEvent {
	return hazards.DepthEvent{Depth: d}
}
func convertByteToDepth(b byte) float64 {
	switch b {
	case 1:
		return 1.0
	case 2:
		return 2.0
	case 3:
		return 3.0
	case 4:
		return 6.0
	case 5:
		return 9.0
	case 7:
		return 0.0 //leveed area
	case 15:
		return 0.0 //intertidal mask only, may experiance high tide or estuarine class in nlcd?
	default:
		return 0.0 //
	}
}
