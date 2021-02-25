package nhc

import (
	"fmt"

	"github.com/USACE/go-consequences/hazards"
	"github.com/dewberry/gdal"
	//"github.com/dewberry/gdal"
)

type nhcInundationData struct {
	FilePath string
}

type Location struct {
	X    float64
	Y    float64
	SRID string
}

type BBox struct {
	bbox []float64
}

func (bb BBox) ToString() string {
	return fmt.Sprintf("%f,%f,%f,%f,%f,%f,%f,%f,%f,%f",
		bb.bbox[0], bb.bbox[1],
		bb.bbox[2], bb.bbox[1],
		bb.bbox[2], bb.bbox[3],
		bb.bbox[0], bb.bbox[3],
		bb.bbox[0], bb.bbox[1])
}

//Init creates and produces an unexported nhcInundationData struct.
func Init(fp string) nhcInundationData {
	//read the file path
	//make sure it is a tif
	return nhcInundationData{FilePath: fp}
}

//ProvideHazard provides a hazardevent for a LocationArgument
func (nid nhcInundationData) ProvideHazard(l Location) (hazards.HazardEvent, error) {
	ds, err := gdal.Open(nid.FilePath, gdal.ReadOnly)
	if err != nil {
		return hazards.DepthEvent{}, err
	}
	defer ds.Close()
	rb := ds.RasterBand(1)
	igt := ds.InvGeoTransform()
	px := int(igt[0] + l.Y*igt[1] + l.X*igt[2])
	py := int(igt[3] + l.Y*igt[4] + l.X*igt[5])
	buffer := make([]float32, 1*1)
	rb.IO(ds.Read, px, py, 1, 1, buffer, 1, 1, 0, 0)
	//return convertDepthtoHazardEvent(convertByteToDepth(3)), nil //buffer[0]))
	return hazards.DepthEvent{convertByteToDepth(3)}, nil
}
func (nid nhcInundationData) GetBoundingBox() (BBox, error) {
	bbox := make([]float64, 4)
	ds, err := gdal.Open(nid.FilePath, gdal.ReadOnly)
	if err != nil {
		return BBox{bbox}, err
	}
	defer ds.Close()

	gt := ds.GeoTransform()
	dx := ds.RasterXSize()
	dy := ds.RasterYSize()
	bbox[0] = gt[0]                     //upper left x
	bbox[1] = gt[2]                     //upper left y
	bbox[2] = gt[0] + gt[1]*float64(dx) //lower right x
	bbox[3] = gt[2] + gt[3]*float64(dy) //lower right y
	return BBox{bbox}, nil
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
