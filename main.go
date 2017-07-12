package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BakeRolls/trends"
	"github.com/wcharczuk/go-chart"
)

func draw(qs []string, iot []trends.IotResult) chart.Chart {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: func(v interface{}) string {
				t := time.Unix(int64(v.(float64)), 0)
				return fmt.Sprintf("%d-%d\n%d", t.Month(), t.Day(), t.Year())
			},
		},
	}

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
			Name:    qs[i],
			XValues: x,
			YValues: y,
		})
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	return graph
}

func svgHandle(res http.ResponseWriter, req *http.Request) {
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
	draw(qs, iot).Render(chart.SVG, res)
}

func csvHandle(res http.ResponseWriter, req *http.Request) {
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
	fmt.Fprint(res, "time")
	for _, q := range qs {
		fmt.Fprintf(res, ",%s", q)
	}
	fmt.Fprintln(res)
	for _, r := range iot {
		fmt.Fprint(res, r.Time.Format("2006-02-01"))
		for _, v := range r.Values {
			fmt.Fprintf(res, ",%d", v)
		}
		fmt.Fprintln(res)
	}
}

func main() {
	host := flag.String("host", ":8080", "Host")
	flag.Parse()

	http.HandleFunc("/", svgHandle)
	http.HandleFunc("/csv", csvHandle)
	http.ListenAndServe(*host, nil)
}
