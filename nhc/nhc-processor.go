package nhc

import (
	"fmt"
	"log"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
	"github.com/dewberry/gdal"
)

type nhcInundationData struct {
	FilePath string
	ds       *gdal.Dataset
}

//Init creates and produces an unexported nhcInundationData struct.
func Init(fp string) nhcInundationData {
	//read the file path
	//make sure it is a tif
	ds, err := gdal.Open(fp, gdal.ReadOnly)
	if err != nil {
		log.Fatalln("Cannot connect to raster.  Killing everything!")
	}
	return nhcInundationData{fp, &ds}
}

func (nid nhcInundationData) Close() {
	nid.ds.Close()
}

//ProvideHazard provides a hazardevent for a LocationArgument
func (nid nhcInundationData) ProvideHazard(l geography.Location) (hazards.HazardEvent, error) {
	rb := nid.ds.RasterBand(1)
	igt := nid.ds.InvGeoTransform()
	px := int(igt[0] + l.X*igt[1] + l.Y*igt[2])
	py := int(igt[3] + l.X*igt[4] + l.Y*igt[5])
	buffer := make([]int32, 1*1)
	rb.IO(gdal.Read, px, py, 1, 1, buffer, 1, 1, 0, 0)
	depth := uint8(buffer[0])
	return convertDepthtoHazardEvent(convertByteToDepth(depth)), nil
}
func (nid nhcInundationData) ProvideHazardBoundary() (geography.BBox, error) {
	bbox := make([]float64, 4)
	gt := nid.ds.GeoTransform()
	fmt.Println(gt)
	dx := nid.ds.RasterXSize()
	dy := nid.ds.RasterYSize()
	bbox[0] = gt[0]                     //upper left x
	bbox[1] = gt[3]                     //upper left y
	bbox[2] = gt[0] + gt[1]*float64(dx) //lower right x
	bbox[3] = gt[3] + gt[5]*float64(dy) //lower right y
	return geography.BBox{Bbox: bbox}, nil
}

func convertDepthtoHazardEvent(d float64) hazards.HazardEvent {
	h := hazards.DepthEvent{}
	h.SetDepth(d)
	return h //could be a hazard.CoastalEvent{Depth:d, Salinity:true}
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
