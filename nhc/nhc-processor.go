package nhc

import (
	"fmt"
	"log"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
	"github.com/dewberry/gdal"
)

type nhcInundationData struct {
	FilePath string
	ds       *gdal.Dataset
}

// Init creates and produces an unexported nhcInundationData struct.
func Init(fp string) nhcInundationData {
	//read the file path
	//make sure it is a tif
	ds, err := gdal.Open(fp, gdal.Access(gdal.ReadOnly))
	if err != nil {
		log.Fatalln("Cannot connect to raster.  Killing everything!")
	}
	return nhcInundationData{fp, &ds}
}

func (nid nhcInundationData) Close() {
	nid.ds.Close()
}

// ProvideHazard provides a hazardevent for a LocationArgument
func (nid nhcInundationData) Hazard(l geography.Location) (hazards.HazardEvent, error) {
	rb := nid.ds.RasterBand(1)
	igt := nid.ds.InvGeoTransform()
	px := int(igt[0] + l.X*igt[1] + l.Y*igt[2])
	py := int(igt[3] + l.X*igt[4] + l.Y*igt[5])
	buffer := make([]int32, 1*1)
	rb.IO(gdal.RWFlag(gdal.Read), px, py, 1, 1, buffer, 1, 1, 0, 0)
	depth := uint8(buffer[0])
	d, err := convertByteToDepthFootBins(depth)
	if err != nil {
		he, heok := err.(hazardproviders.HazardError)
		if heok {
			qe := hazards.QualitativeEvent{}
			qe.SetQualitative(he.Input)
			return qe, nil
		}
		return hazards.DepthEvent{}, err
	}
	return convertDepthtoHazardEvent(d), nil
}
func (nid nhcInundationData) HazardBoundary() (geography.BBox, error) {
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
func convertByteToDepth(b byte) (float64, error) {
	switch b {
	case 1:
		return 1.0, nil
	case 2:
		return 2.0, nil
	case 3:
		return 3.0, nil
	case 4:
		return 6.0, nil
	case 5:
		return 9.0, nil
	case 7:
		return -901.0, hazardproviders.HazardError{"Leveed Area detected"} //leveed area
	case 15:
		return -901.0, hazardproviders.HazardError{"Inter Tidal Mask detected"} //intertidal mask only, may experiance high tide or estuarine class in nlcd?
	default:
		return -901.0, hazardproviders.NoHazardFoundError{"Byte value of " + string(b) + "is not tracked as a hazard."} //
	}
}
func convertByteToDepthFootBins(b byte) (float64, error) {
	switch b {
	case 1:
		return 1.0, nil
	case 2:
		return 2.0, nil
	case 3:
		return 3.0, nil
	case 4:
		return 4.0, nil
	case 5:
		return 5.0, nil
	case 6:
		return 6.0, nil
	case 7:
		return 7.0, nil
	case 8:
		return 8.0, nil
	case 9:
		return 9.0, nil
	case 10:
		return 10.0, nil
	case 11:
		return 11.0, nil
	case 12:
		return 12.0, nil
	case 13:
		return 13.0, nil
	case 14:
		return 14.0, nil
	case 15:
		return 15.0, nil
	case 16:
		return 16.0, nil
	case 17:
		return 17.0, nil
	case 18:
		return 18.0, nil
	case 19:
		return 19.0, nil
	case 20:
		return 20.0, nil
	case 21:
		return 21.0, nil
	case 99:
		return -901.0, hazardproviders.HazardError{"Leveed Area detected"} //leveed area
	case 88:
		return -901.0, hazardproviders.HazardError{"Inter Tidal Mask detected"} //intertidal mask only, may experiance high tide or estuarine class in nlcd?
	default:
		return -901.0, hazardproviders.NoHazardFoundError{"Byte value of " + string(b) + "is not tracked as a hazard."} //
	}
}
