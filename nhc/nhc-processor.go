package nhc

import (
	"github.com/USACE/go-consequences/hazards"
	//"github.com/dewberry/gdal"
)

type nhcInundationData struct {
	FilePath string
}
type LocationArgument struct {
	X float64
	Y float64
}

//Init creates and produces an unexported nhcInundationData struct.
func Init(fp string) nhcInundationData {
	//read the file path
	//make sure it is a tif
	return nhcInundationData{FilePath: fp}
}

//ProvideHazard provides a hazardevent for a LocationArgument
func (nid nhcInundationData) ProvideHazard(args interface{}) (hazards.HazardEvent, error) {
	//needs work.
	//la, ok := args.(LocationArgument)
	//if ok{
	/*x := la.X
	y := la.Y
	ds, err := gdal.Open(nid.FilePath, gdal.ReadOnly)
	if err != nil {
		return 0.0, err
	}
	rb := ds.RasterBand(1)
	igt := ds.InvGeoTransform()
	px := int(igt[0] + y*igt[1] + x*igt[2])
	py := int(igt[3] + y*igt[4] + x*igt[5])
	buffer := make([]float32, 1*1)
	rb.IO(gdal.Read, px, py, 1, 1, buffer, 1, 1, 0, 0)
	*/
	return convertDepthtoHazardEvent(convertByteToDepth(3)), nil //buffer[0]))
	//}
	//err := hazardproviders.HazardError{Input: "Could not Parse args"}
	//return nil, err
}
func (nid nhcInundationData) GetBoundingBox() string {
	//needs to be in format NSI expects
	return "-81.58418,30.25165,-81.58161,30.26939,-81.55898,30.26939,-81.55281,30.24998,-81.58418,30.25165"
}
func convertDepthtoHazardEvent(d float64) hazards.HazardEvent {
	return hazards.DepthEvent{Depth: d} //could be a hazard.CoastalEvent{Depth:d, Salinity:true}
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
