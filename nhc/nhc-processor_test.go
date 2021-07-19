package nhc

import (
	"fmt"
	"testing"

	"github.com/HydrologicEngineeringCenter/go-tc-consequences/outputwriter"
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structureprovider"
)

func Test_Convert(t *testing.T) {
	b, err := convertByteToDepth(1)
	if err != nil {
		panic(err)
	}
	d := convertDepthtoHazardEvent(b)
	de, _ := d.(hazards.DepthEvent)
	if de.Depth() != 1.0 {
		t.Errorf("expected 1.0 got something else.")
	}
	b2, err := convertByteToDepth(2)
	if err != nil {
		panic(err)
	}
	d2 := convertDepthtoHazardEvent(b2)
	de2, _ := d2.(hazards.DepthEvent)
	if de2.Depth() != 2.0 {
		t.Errorf("expected 2.0 got something else.")
	}
	b3, err := convertByteToDepth(3)
	if err != nil {
		panic(err)
	}
	d3 := convertDepthtoHazardEvent(b3)
	de3, _ := d3.(hazards.DepthEvent)
	if de3.Depth() != 3.0 {
		t.Errorf("expected 3.0 got something else.")
	}
	b4, err := convertByteToDepth(4)
	if err != nil {
		panic(err)
	}
	d4 := convertDepthtoHazardEvent(b4)
	de4, _ := d4.(hazards.DepthEvent)
	if de4.Depth() != 6.0 {
		t.Errorf("expected 6.0 got something else.")
	}
	b5, err := convertByteToDepth(5)
	if err != nil {
		panic(err)
	}
	d5 := convertDepthtoHazardEvent(b5)
	de5, _ := d5.(hazards.DepthEvent)
	if de5.Depth() != 9.0 {
		t.Errorf("expected 9.0 got something else.")
	}
	b6, err := convertByteToDepth(6)
	if err == nil {
		t.Errorf("expected error about byte 6 not being acceptable got something else.")
	}
	fmt.Printf("%f\n", b6)
	b7, err := convertByteToDepth(7)
	if err == nil {
		t.Errorf("expected error about levees got something else.")
	}
	fmt.Printf("%f\n", b7)
	b15, err := convertByteToDepth(15)
	if err != nil {
		panic(err)
	}
	d15 := convertDepthtoHazardEvent(b15)
	de15, _ := d15.(hazards.DepthEvent)
	if de15.Depth() != 0.0 {
		t.Errorf("expected 0.0 got something else.")
	}
}
func Test_Compute_shp(t *testing.T) {
	hp := Init("/workspaces/go-tc-consequences/data/LakeC_LAURA_2020_adv19_e10_ResultMaskRaster_4326.tif")
	sp, se := structureprovider.InitGPK("/workspaces/go-tc-consequences/data/nsi.gpkg", "nsi")
	if se != nil {
		panic(se)
	}
	ow, err := outputwriter.InitNHCShpResultsWriter("/workspaces/go-tc-consequences/data/LakeC_LAURA_2020_adv19_e10_ResultMaskRaster.shp", "NHC_RESULTS")
	if err != nil {
		panic(err)
	}
	defer ow.Close()
	compute.StreamAbstract(hp, sp, ow)
}
func Test_Compute_gpkg(t *testing.T) {
	hp := Init("/workspaces/go-tc-consequences/data/LakeC_LAURA_2020_adv19_e10_ResultMaskRaster_4326.tif")
	sp, se := structureprovider.InitGPK("/workspaces/go-tc-consequences/data/nsi.gpkg", "nsi")
	if se != nil {
		panic(se)
	}
	ow, err := outputwriter.InitNHCGpkResultsWriter("/workspaces/go-tc-consequences/data/LakeC_LAURA_2020_adv19_e10_ResultMaskRaster.gpkg", "NHC_RESULTS")
	if err != nil {
		panic(err)
	}
	defer ow.Close()
	compute.StreamAbstract(hp, sp, ow)
}
func Test_Compute_json(t *testing.T) {
	hp := Init("/workspaces/go-tc-consequences/data/clipped_sample.tif")
	sp, se := structureprovider.InitGPK("/workspaces/go-tc-consequences/data/nsi.gpkg", "nsi")
	if se != nil {
		panic(se)
	}
	ow, err := consequences.InitGeoJsonResultsWriterFromFile("/workspaces/go-tc-consequences/data/clipped_sample_consequences.json")
	if err != nil {
		panic(err)
	}
	defer ow.Close()
	compute.StreamAbstract(hp, sp, ow)
}
