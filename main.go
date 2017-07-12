package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/BakeRolls/trends"
	"github.com/wcharczuk/go-chart"
)

func draw(iot []trends.IotResult) chart.Chart {
	graph := chart.Chart{}

	for i := 0; i < len(iot[0].Values); i++ {
		x := []float64{}
		for _, r := range iot {
			x = append(x, float64(r.Time.Unix()))
		}

		y := []float64{}
		for _, r := range iot {
			y = append(y, float64(r.Values[i]))
		}

		graph.Series = append(graph.Series, chart.ContinuousSeries{
			XValues: x,
			YValues: y,
		})
	}

	return graph
}

func index(res http.ResponseWriter, req *http.Request) {
	q, ok := req.URL.Query()["q"]
	if !ok || len(q) == 0 {
		fmt.Fprintln(res, "?q=foo,bar")
		return
	}
	log.Println(q)
	qs := strings.Split(q[0], ",")

	iot, err := trends.InterestOverTime(qs...)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	if len(iot) == 0 {
		fmt.Fprint(res, "unexpected response")
		return
	}

	res.Header().Set("Content-Type", chart.ContentTypeSVG)
	draw(iot).Render(chart.SVG, res)
}

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}
