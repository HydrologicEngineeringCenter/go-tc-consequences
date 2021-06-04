package main

import (
	"github.com/HydrologicEngineeringCenter/go-tc-consequences/nhc"
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/structureprovider"
)

func main() {
	//serverless solution
	argsWithoutProg := os.Args[1:]

	//-ss (structures source) -sfp (structure file path)  -hs (hazard source) -hfp (hazard file path) -ot (output type) //we will define the path internally?
	if len(argsWithoutProg) != 2 {
		fmt.Println("Expected two arguments, the filepath to the csv input and the file path to the geopackage input")
	} else {
		hfp := argsWithoutProg[0]
		sfp := argsWithoutProg[1]
		fmt.Println(fmt.Sprintf("Computing EAD for %v using an iventory at path %v", hfp, sfp))
		compute.ExpectedAnnualDamagesGPK(hfp, sfp)
	}

	root := "/workspaces/go-tc-consequences/data/clipped_sample"
	filepath := root + ".tif"
	jwriter := consequences.InitJsonResultsWriterFromFile(root + "_consequences.json")
	nsp := structureprovider.InitNSISP()
	defer jwriter.Close()
	nhcTiffReader := nhc.Init(filepath)

	compute.StreamAbstract(nhcTiffReader, nsp, jwriter)

}
