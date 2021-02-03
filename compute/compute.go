package compute

import (
	"github.com/USACE/go-consequences/nsi"
)

func compute(filepath string) {
	//read the file path
	//make sure it is a tif
	//open a tif reader
	//get boundingbox

	//
	nsi.GetByBboxStream("convertboundingboxtostring", func(f nsi.NsiFeature) {
		//convert nsifeature to structure
		//query input tiff for xy location
		//convert bytes to depths.
		//compute damages based on provided depths
		//keep a summmary of damages
	})
}
