package main

import (
	"github.com/HydrologicEngineeringCenter/go-tc-consequences/nhc"
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/structureprovider"
)

func main() {
	//serverless solution
	root := "/workspaces/go-tc-consequences/data/clipped_sample"
	filepath := root + ".tif"
	jwriter := consequences.InitJsonResultsWriterFromFile(root + "_consequences.json")
	nsp := structureprovider.InitNSISP()
	defer jwriter.Close()
	nhcTiffReader := nhc.Init(filepath)

	compute.StreamAbstract(nhcTiffReader, nsp, jwriter)
	/*//server solution
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		fp, fpPresent := params["FilePath"]
		if !fpPresent {
			http.Error(w, "No FilePath argument", http.StatusNotFound)
		} else {
			if len(fp[0]) == 0 {
				//should have better error checking...
				http.Error(w, "Invalid FilePath argument", http.StatusNotFound)
			} else {
				//fmt.Fprintf(w, fp[0])
				compute.ComputeFromFilePathWithWriter(fp[0], w)
			}
		}
	})
	log.Print("starting local server")
	log.Fatal(http.ListenAndServe("localhost:3030", nil))
	*/

}
