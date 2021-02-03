package nhc

import (
	"github.com/USACE/go-consequences/hazards"
	//"github.com/dewberry/gdal"
)

type nhcInundationData struct {
	FilePath string
}

//Init creates and produces an unexported nhcInundationData struct.
func Init(fp string) nhcInundationData {
	//read the file path
	//make sure it is a tif
	return nhcInundationData{FilePath: fp}
}

//GetHazardEvent should get converted to go-consequences hazardprovider.ProvideHazard(args interface{})hazards.HazardEvent,err
func (nid nhcInundationData) GetHazardEvent(x float64, y float64) hazards.HazardEvent {
	//needs work.
	/*ds, err := gdal.Open(nid.FilePath, gdal.ReadOnly)
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
	return convertDepthtoDepthEvent(convertByteToDepth(1)) //buffer[0]))
}
func (nid nhcInundationData) GetBoundingBox() string {
	//needs to be in format NSI expects
	return "gobbldygook"
}
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
