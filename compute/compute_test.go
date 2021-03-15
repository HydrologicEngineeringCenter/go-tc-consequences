package compute

import (
	"testing"

	"github.com/HenryGeorgist/go-statistics/statistics"
	"github.com/HydrologicEngineeringCenter/go-tc-consequences/nhc"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structures"
)

func Test_Compute_NSI(t *testing.T) {
	t.Log("Starting Test")
	//compute("/vsis3/usace-storms/cwbi-dls-mmc/tiles/FLOOD_DEPTH/2307_COG.tif")
	ComputeFromFilePath("/workspaces/go-tc-consequences/data/clipped_sample.tif")
}
func Test_Compute_FakeHP(t *testing.T) {
	mhp := mockhp{}
	t.Log("Starting Mockupped Test")
	compute(mhp)
}
func Test_Compute_FakeStructure(t *testing.T) {
	//get a map of all occupancy types
	m := structures.OccupancyTypeMap()
	//pick one for testing
	var o = m["RES1-1SNB"]
	//define a distribution of strucure value
	sv := statistics.NormalDistribution{Mean: 100.00, StandardDeviation: 1}
	//define a distribution of content value
	cv := statistics.NormalDistribution{Mean: 100.00, StandardDeviation: 1}
	//mutate to a ParameterValue for homogeneity
	spv := consequences.ParameterValue{Value: sv}
	cpv := consequences.ParameterValue{Value: cv}
	//foundation height of zero (as a constant rather than a distribution)
	fhpv := consequences.ParameterValue{Value: 0}
	//create a fake structure
	var s = structures.StructureStochastic{OccType: o, StructVal: spv, ContVal: cpv, FoundHt: fhpv, BaseStructure: structures.BaseStructure{DamCat: "category", X: 1.0, Y: 1.0}}
	//turn uncertainty off to test with mean values.
	s.UseUncertainty = false
	//get unexported tiff reader.
	tiffReader := nhc.Init("ultimatelythisneedstobeatif")
	//query input tiff for xy location (fake x,y for testing)
	d, _ := tiffReader.ProvideHazard(geography.Location{X: s.BaseStructure.X, Y: s.BaseStructure.Y})
	//cast to float - because it is an empty interface?
	got := s.Compute(d).Result.Result[0].(float64) //zero is structure damage not content damage.
	if got-40.099998 > .0000005 {
		t.Errorf("Compute() = %f", got)
	}

}

//developing mockups for testing
type mockhp struct {
}

//ProvideHazard provides a hazardevent for a LocationArgument
func (hp mockhp) ProvideHazard(l geography.Location) (hazards.HazardEvent, error) {
	h := hazards.DepthEvent{}
	h.SetDepth(3.0)
	return h, nil
}
func (hp mockhp) ProvideHazardBoundary() (geography.BBox, error) {
	bbox := make([]float64, 4) //i might have these values inverted
	bbox[0] = -81.58418        //upper left x
	bbox[1] = 30.25165         //upper left y
	bbox[2] = -81.58161        //lower right x
	bbox[3] = 30.26939         //lower right y
	return geography.BBox{Bbox: bbox}, nil
}
func (hp mockhp) Close() {
	//do nothing.
}
