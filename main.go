package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/BakeRolls/test-http/handler"
	"github.com/BakeRolls/trends"
	"github.com/justinas/alice"
	"github.com/wcharczuk/go-chart"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<ul>
			<li><a href="/csv?q=foo&q=bar">/csv</a></li>
			<li><a href="/png?q=foo&q=bar">/png</a></li>
			<li><a href="/svg?q=foo&q=bar">/svg</a></li>
			<li><a href="https://github.com/BakeRolls/trends-http">github</a></li>
		</ul>
	`)
}

func csv(w http.ResponseWriter, r *http.Request) {
	names := r.Context().Value(handler.NamesKey).([]string)
	iot := r.Context().Value(handler.IotKey).([]trends.IotResult)

	fmt.Fprint(w, "time")
	for _, n := range names {
		fmt.Fprintf(w, ",%s", n)
	}
	fmt.Fprintln(w)
	for _, r := range iot {
		fmt.Fprint(w, r.Time.Format("2006-02-01"))
		for _, v := range r.Values {
			fmt.Fprintf(w, ",%d", v)
		}
		fmt.Fprintln(w)
	}
}

func png(w http.ResponseWriter, r *http.Request) {
	names := r.Context().Value(handler.NamesKey).([]string)
	iot := r.Context().Value(handler.IotKey).([]trends.IotResult)

	w.Header().Set("Content-Type", chart.ContentTypePNG)
	draw(names, iot).Render(chart.PNG, w)
}

func svg(w http.ResponseWriter, r *http.Request) {
	names := r.Context().Value(handler.NamesKey).([]string)
	iot := r.Context().Value(handler.IotKey).([]trends.IotResult)

	w.Header().Set("Content-Type", chart.ContentTypeSVG)
	draw(names, iot).Render(chart.SVG, w)
}

func main() {
	host := flag.String("host", ":8080", "Host")
	flag.Parse()

	handlers := alice.New(handler.Init, handler.Log)
	iotHandlers := handlers.Append(handler.Iot)

	http.Handle("/", handlers.ThenFunc(index))
	http.Handle("/csv", iotHandlers.ThenFunc(csv))
	http.Handle("/png", iotHandlers.ThenFunc(png))
	http.Handle("/svg", iotHandlers.ThenFunc(svg))
	http.ListenAndServe(*host, nil)
}
