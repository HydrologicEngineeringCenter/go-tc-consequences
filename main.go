package main

import (
	"os"

	"github.com/HydrologicEngineeringCenter/go-tc-consequences/compute"
)

func main() {
	//serverless solution
	filepath := "/workspaces/go-tc-consequences/data/clipped_sample.tif"
	w, err := os.OpenFile("/workspaces/go-tc-consequences/data/clipped_sample_consequences.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	compute.ComputeFromFilePathWithWriter(filepath, w)
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
