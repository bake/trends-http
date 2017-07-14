package main

import (
	"fmt"
	"time"

	"github.com/BakeRolls/trends"
	"github.com/wcharczuk/go-chart"
)

func draw(names []string, iot []trends.IotResult) chart.Chart {
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
			Name:    names[i],
			XValues: x,
			YValues: y,
		})
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	return graph
}
