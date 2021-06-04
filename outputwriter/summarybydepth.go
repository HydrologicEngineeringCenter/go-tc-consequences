package outputwriter


import (
	"fmt"
	"io"
	"os"
)

type summaryByDepth struct {
	filepath   string
	w          io.Writer
	grandTotal int
	thresholds []int{0,2,4,6}
	headers	[]string{"No Damage (0 ft)", "Affected (<=2 ft)", "Minor Damage (2 - 4 ft)", "Major Damage (4 - 6 ft)", "Destroyed (6+ ft)"}
	totals     []int
}

func InitSummaryByDepth(filepath string) *summaryByDepth {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	t := make([]int, 5)
	return &summaryByDepth{filepath: filepath, w: w, totals: t}
}
func (srw *summaryByDepth) Write(r Result) {
	//hardcoding for structures to experiment and think it through.
	value :=0.0
	for i, val := range r.Headers {
		if val == "hazard"{
				de, dok := value.(hazards.HazardEvent)
				if dok {
					if de.Has(hazards.Depth) {
						value = de.Depth()
						
					}
				} else {
					//must be an array - bummer.
					//get at the elements of the slice, add all depths to the table?
				}
		}
	}
	counted := false
	for i, val := range srw.thresholds{
		if value <= val{
			srw.totals[i] +=1
			counted = true
			srw.grandTotal +=1
			break
		}
	}
	if !counted{
		if value >0 {
			srw.totals[len(srw.totals)-1] += 1
			srw.grandTotal +=1
		}
	}
}
func (srw *summaryByDepth) Close() {
	fmt.Fprintf(srw.w, "Outcome, Count\n")
	h := srw.totals
	for i, v := range h {
		fmt.Fprintf(srw.w, fmt.Sprintf("%v, %v\n", srw.headers[i], v))
	}
	fmt.Fprintf(srw.w, fmt.Sprintf("Total Building Count, %v\n", srw.grandTotal))
	w2, ok := srw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}
