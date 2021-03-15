package compute

import (
	"fmt"
	"io"
	"log"

	"github.com/HydrologicEngineeringCenter/go-tc-consequences/nhc"
	comp "github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/nsi"
	"github.com/USACE/go-consequences/structures"
)

func ComputeFromFilePath(filepath string) {
	tiffReader := nhc.Init(filepath)
	defer tiffReader.Close()
	compute(tiffReader)
}
func compute(hp hazardproviders.HazardProvider) {

	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
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
		d, _ := hp.ProvideHazard(geography.Location{X: str.X, Y: str.Y})
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
func ComputeFromFilePathWithWriter(filepath string, w io.Writer) {
	tiffReader := nhc.Init(filepath)
	defer tiffReader.Close()
	computeWithWriter(tiffReader, w)
}
func computeWithWriter(hp hazardproviders.HazardProvider, w io.Writer) {
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	//get a map of all occupancy types
	m := structures.OccupancyTypeMap()
	//define a default occtype in case of emergancy
	defaultOcctype := m["RES1-1SNB"]
	//create a header for marshalling
	header := []string{"fd_id", "x", "y", "depth", "structure damage", "content damage", "Pop_2amo65", "Pop_2amu65", "Pop_2pmo65", "Pop_2pmu65"}

	nsi.GetByBboxStream(bbox.ToString(), func(f nsi.NsiFeature) {
		//convert nsifeature to structure
		str := comp.NsiFeaturetoStructure(f, m, defaultOcctype)
		//query input tiff for xy location
		d, _ := hp.ProvideHazard(geography.Location{X: str.X, Y: str.Y})
		//compute damages based on provided depths
		if d.Has(hazards.Depth) {
			if d.Depth() > 0.0 {
				r := str.Compute(d)
				//keep a summmary of damages that adds the structure name
				row := []interface{}{r.Result.Result[0], r.Result.Result[1], r.Result.Result[2], r.Result.Result[3], r.Result.Result[4], r.Result.Result[5], f.Properties.Pop2amo65, f.Properties.Pop2amu65, f.Properties.Pop2pmo65, f.Properties.Pop2pmu65}
				structureResult := consequences.Result{Headers: header, Result: row}
				b, _ := structureResult.MarshalJSON()
				s := string(b) + "\n"
				fmt.Fprintf(w, s)
			}
		}
	})
}
