package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/HydrologicEngineeringCenter/go-tc-consequences/nhc"
	"github.com/HydrologicEngineeringCenter/go-tc-consequences/outputwriter"
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
)

func main() {
	//command line program
	var err error
	var sp consequences.StreamProvider
	var hp hazardproviders.HazardProvider //not sure what it will be yet, but we can declare it!
	hazardProviderVerticalUnitsIsFeet := true
	var ow consequences.ResultsWriter //need a file path to write anything...

	//-ss (structures source) -sfp (structure file path)  -hs (hazard source) -hfp (hazard file path) -ot (output type) //we will define the path internally?
	sfp := flag.String("sfp", "", "structure file path, (optional)")
	ss := flag.String("ss", "nsi", "structure source, (optional), acceptable terms: nsi(default), gpkg, shp")
	hfp := flag.String("hfp", "", "hazard file path, (required)")
	hs := flag.String("hs", "", "hazard source, (required), acceptable terms: nhc, depths")
	ht := flag.String("ht", "feet", "hazard type of vertical datum, (required), acceptable terms: feet (default), meters")
	ot := flag.String("ot", "gpkg", "output type, (optional), acceptable terms: gpkg (default), shp, geojson, summaryDollars, summaryDepths")
	var se error
	se = nil
	flag.Parse()
	if *sfp != "" {
		switch *ss {
		case "gpkg":
			sp, se = structureprovider.InitGPK(*sfp, "nsi")

		case "shp":
			sp, se = structureprovider.InitSHP(*sfp)
		case "nsi":
			sp = structureprovider.InitNSISP() //default to NSI API structure provider.
		default:
			sp = structureprovider.InitNSISP()
		}
	} else {
		sp = structureprovider.InitNSISP()
	}
	if *ht != "" {
		switch *ht {
		case "feet":
			hazardProviderVerticalUnitsIsFeet = true
		case "meters":
			hazardProviderVerticalUnitsIsFeet = false
		}
	}
	var he error
	he = nil
	if *hfp != "" {
		switch *hs {
		case "nhc":
			hp = nhc.Init(*hfp)
		case "depths":
			if hazardProviderVerticalUnitsIsFeet {
				hp, he = hazardproviders.Init(*hfp)
			} else {
				hp, he = hazardproviders.Init_Meters(*hfp)
			}
		}
	} else {
		he = errors.New("cannot compute without hazard provider path, use -h for help.")
	}
	ofp := *hfp
	// pull the .tif off the end?
	ofp = ofp[:len(ofp)-4] //good enough for government work?
	fmt.Println(ofp)
	var oe error
	oe = nil
	if ofp != "" {
		switch *ot {
		case "gpkg":
			ofp += "_consequences.gpkg"
			if *hs == "nhc" {
				ow, oe = outputwriter.InitNHCGpkResultsWriter(ofp, "results")
			} else {
				ow, oe = resultswriters.InitGpkResultsWriter(ofp, "results")
			}

		case "shp":
			ofp += "_consequences.shp"
			if *hs == "nhc" {
				ow, oe = outputwriter.InitNHCShpResultsWriter(ofp, "results")
			} else {
				ow, oe = resultswriters.InitShpResultsWriter(ofp, "results")
			}

		case "geojson":
			ofp += "_consequences.json"
			ow, oe = resultswriters.InitGeoJsonResultsWriterFromFile(ofp)
		case "summaryDollars":
			ofp += "_summaryDollars.csv"
			ow, oe = resultswriters.InitSummaryResultsWriterFromFile(ofp)
		case "summaryDepths":
			ofp += "_summaryDepths.csv"
			ow = outputwriter.InitSummaryByDepth(ofp)
		default:
			ofp += "_consequences.gpkg"
			ow, oe = resultswriters.InitGpkResultsWriter(ofp, "results")
		}
	} else {
		oe = errors.New("we need an input hazard file path use -h for help.")
	}
	defer ow.Close()
	if se != nil {
		if he != nil {
			if oe != nil {
				err = errors.New(se.Error() + "\n" + he.Error() + "\n" + oe.Error() + "\n")
			} else {
				err = errors.New(se.Error() + "\n" + he.Error() + "\n")
			}
		} else {
			err = errors.New(se.Error() + "\n")
		}
	} else if he != nil {
		if oe != nil {
			err = errors.New(he.Error() + "\n" + oe.Error() + "\n")
		} else {
			err = errors.New(he.Error() + "\n")
		}
	} else {
		if oe != nil {
			err = errors.New(oe.Error() + "\n")
		}
	}
	if err != nil {
		fmt.Println(err)
	} else {
		compute.StreamAbstract(hp, sp, ow)
	}

}
