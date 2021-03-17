package compute

import (
	"fmt"
	"io"
	"log"

	"github.com/HydrologicEngineeringCenter/go-tc-consequences/nhc"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structureprovider"
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
	header := []string{"fd_id", "x", "y", "depth", "damage_category", "occupancy_type", "structure damage", "content damage", "Pop_2amo65", "Pop_2amu65", "Pop_2pmo65", "Pop_2pmu65"}

	var rows []interface{}
	result := consequences.Results{IsTable: true}
	result.Result.Headers = header
	result.Result.Result = rows
	structureprovider.GetByBboxStream(bbox.ToString(), func(f structureprovider.NsiFeature) {
		//convert nsifeature to structure
		str := structureprovider.NsiFeaturetoStructure(f, m, defaultOcctype)
		//query input tiff for xy location
		d, _ := hp.ProvideHazard(geography.Location{X: str.X, Y: str.Y})
		//compute damages based on provided depths
		if d.Has(hazards.Depth) {
			if d.Depth() > 0.0 {
				r := str.Compute(d)
				//keep a summmary of damages that adds the structure name
				row := []interface{}{r.Result[0], r.Result[1], r.Result[2], d.Depth(), r.Result[4], r.Result[5], r.Result[6], r.Result[7], f.Properties.Pop2amo65, f.Properties.Pop2amu65, f.Properties.Pop2pmo65, f.Properties.Pop2pmu65}
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
	header := []string{"fd_id", "x", "y", "depth", "damage_category", "occupancy_type", "structure damage", "content damage", "Pop_2amo65", "Pop_2amu65", "Pop_2pmo65", "Pop_2pmu65"}

	structureprovider.GetByBboxStream(bbox.ToString(), func(f structureprovider.NsiFeature) {
		//convert nsifeature to structure
		str := structureprovider.NsiFeaturetoStructure(f, m, defaultOcctype)
		//query input tiff for xy location
		d, _ := hp.ProvideHazard(geography.Location{X: str.X, Y: str.Y})
		//compute damages based on provided depths
		if d.Has(hazards.Depth) {
			if d.Depth() > 0.0 {
				r := str.Compute(d)
				//keep a summmary of damages that adds the structure name
				row := []interface{}{r.Result[0], r.Result[1], r.Result[2], d.Depth(), r.Result[4], r.Result[5], r.Result[6], r.Result[7], f.Properties.Pop2amo65, f.Properties.Pop2amu65, f.Properties.Pop2pmo65, f.Properties.Pop2pmu65}
				structureResult := consequences.Result{Headers: header, Result: row}
				b, _ := structureResult.MarshalJSON()
				s := string(b) + "\n"
				fmt.Fprintf(w, s)
			}
		}
	})
}
