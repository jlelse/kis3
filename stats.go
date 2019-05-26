package main

import (
	"encoding/json"
	"fmt"
	"github.com/kis3/kis3/helpers"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func initStatsRouter()  {
	app.router.HandleFunc("/stats", StatsHandler)
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	// Require authentication
	if appConfig.statsAuth() {
		if !helpers.CheckAuth(w, r, appConfig.StatsUsername, appConfig.StatsPassword) {
			return
		}
	}
	// Do request
	queries := r.URL.Query()
	view := PAGES
	switch strings.ToLower(queries.Get("view")) {
	case "pages":
		view = PAGES
	case "referrers":
		view = REFERRERS
	case "useragents":
		view = USERAGENTS
	case "useragentnames":
		view = USERAGENTNAMES
	case "hours":
		view = HOURS
	case "days":
		view = DAYS
	case "weeks":
		view = WEEKS
	case "months":
		view = MONTHS
	case "allhours":
		view = ALLHOURS
	case "alldays":
		view = ALLDAYS
	case "count":
		view = COUNT
	}
	result, e := request(&ViewsRequest{
		view:     view,
		from:     queries.Get("from"),
		fromRel:  queries.Get("fromrel"),
		to:       queries.Get("to"),
		toRel:    queries.Get("torel"),
		url:      queries.Get("url"),
		ref:      queries.Get("ref"),
		ua:       queries.Get("ua"),
		ordercol: strings.ToLower(queries.Get("ordercol")),
		order:    strings.ToUpper(queries.Get("order")),
		limit:    queries.Get("limit"),
	})
	if e != nil {
		fmt.Println("Database request failed:", e)
		w.WriteHeader(500)
	} else if result != nil {
		w.Header().Set("Cache-Control", "max-age=0")
		switch queries.Get("format") {
		case "json":
			sendJsonResponse(result, w)
		case "chart":
			sendChartResponse(result, w)
		default: // "plain"
			sendPlainResponse(result, w)
		}
	}
}

func sendPlainResponse(result []*RequestResultRow, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	for _, row := range result {
		_, _ = fmt.Fprintln(w, (*row).First+": "+strconv.Itoa((*row).Second))
	}
}

func sendJsonResponse(result []*RequestResultRow, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	jsonBytes, _ := json.Marshal(result)
	_, _ = fmt.Fprintln(w, string(jsonBytes))
}

func sendChartResponse(result []*RequestResultRow, w http.ResponseWriter) {
	labels := make([]string, len(result))
	values := make([]int, len(result))
	for i, row := range result {
		labels[i] = row.First
		values[i] = row.Second
	}
	chartJSString, e := app.staticBox.FindString("Chart.min.js")
	if e != nil {
		return
	}
	data := struct {
		Labels  []string
		Values  []int
		ChartJS template.JS
	}{
		Labels:  labels,
		Values:  values,
		ChartJS: template.JS(chartJSString),
	}
	chartTemplateString, e := app.staticBox.FindString("chart.html")
	if e != nil {
		return
	}
	t, e := template.New("chart").Parse(chartTemplateString)
	if e != nil {
		return
	}
	_ = t.Execute(w, data)
}
