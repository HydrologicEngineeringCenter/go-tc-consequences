package compute

import (
	"fmt"

	"github.com/HydrologicEngineeringCenter/go-tc-consequences/nhc"
	comp "github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/nsi"
	"github.com/USACE/go-consequences/structures"
)

func compute(filepath string) {
	//open a tif reader
	tiffReader := nhc.Init(filepath)
	//get boundingbox
	bbox := tiffReader.GetBoundingBox()
	//get a map of all occupancy types
	m := structures.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	nsi.GetByBboxStream(bbox, func(f nsi.NsiFeature) {
		//convert nsifeature to structure
		str := comp.NsiFeaturetoStructure(f, m, defaultOcctype)
		//query input tiff for xy location
		d := tiffReader.GetHazardEvent(str.X, str.Y)
		//compute damages based on provided depths
		de, ok := d.(hazards.DepthEvent)
		if ok {
			if de.Depth > 0.0 {
				r := str.Compute(de)
				//keep a summmary of damages
				fmt.Println(r)
			}
		}
	})
}
