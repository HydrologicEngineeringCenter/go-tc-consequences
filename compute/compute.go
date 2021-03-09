package compute

import (
	"fmt"
	"log"

	"github.com/HydrologicEngineeringCenter/go-tc-consequences/nhc"
	comp "github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/nsi"
	"github.com/USACE/go-consequences/structures"
)

func compute(filepath string) {
	//open a tif reader
	tiffReader := nhc.Init(filepath)
	defer tiffReader.Close()
	//get boundingbox
	fmt.Println("Getting bbox")
	bbox, err := tiffReader.GetBoundingBox()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	//get a map of all occupancy types
	m := structures.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	//create a results store
	header := []string{"fd_id", "x", "y", "structure damage", "content damage"}
	var rows []interface{}
	result := consequences.Results{IsTable: true}
	result.Result.Headers = header
	result.Result.Result = rows
	nsi.GetByBboxStream(bbox.ToString(), func(f nsi.NsiFeature) {
		//convert nsifeature to structure
		str := comp.NsiFeaturetoStructure(f, m, defaultOcctype)
		//query input tiff for xy location
		d, _ := tiffReader.ProvideHazard(nhc.Location{X: str.X, Y: str.Y})
		//compute damages based on provided depths
		if d.Has(hazards.Depth) {
			if d.Depth() > 0.0 {
				r := str.Compute(d)
				//keep a summmary of damages that adds the structure name
				row := []interface{}{str.Name, str.X, str.Y, r.Result.Result[0], r.Result.Result[1]}
				structureResult := consequences.Result{Headers: header, Result: row}
				result.AddResult(structureResult)
			}
		}
	})
	b, _ := result.MarshalJSON() //json.Marshal(result)
	fmt.Println(string(b))
	//fmt.Println(result)
}
