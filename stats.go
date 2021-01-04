package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"git.jlel.se/jlelse/kis3/helpers"
)

func initStatsRouter() {
	app.router.HandleFunc("/stats", statsHandler)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	// Require authentication
	if appConfig.statsAuth() {
		if !helpers.CheckAuth(w, r, appConfig.StatsUsername, appConfig.StatsPassword) {
			return
		}
	}
	// Do request
	queryValues := r.URL.Query()
	result, e := doRequest(queryValues)
	if e != nil {
		fmt.Println("Database request failed:", e)
		w.WriteHeader(500)
	} else if result != nil {
		w.Header().Set("Cache-Control", "max-age=0")
		switch queryValues.Get("format") {
		case "json":
			sendJSONResponse(result, w)
		case "chart":
			sendChartResponse(result, w)
		default: // "plain"
			sendPlainResponse(result, w)
		}
	}
}

func doRequest(queries url.Values) (result []*RequestResultRow, e error) {
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
	case "os":
		view = OS
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
	result, e = request(&ViewsRequest{
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
		bots:     queries.Get("bots"),
		os:       queries.Get("os"),
	})
	return
}

func sendPlainResponse(result []*RequestResultRow, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	for _, row := range result {
		_, _ = fmt.Fprintln(w, (*row).First+": "+strconv.Itoa((*row).Second))
	}
}

func sendJSONResponse(result []*RequestResultRow, w http.ResponseWriter) {
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

func respondToTelegramUpdate(u *telegramUpdate) {
	if app.telegram != nil && strings.HasPrefix(u.Message.Text, "/stats") {
		response := ""
		fakeURL, e := url.Parse("/stats?" + strings.TrimSpace(strings.TrimPrefix(u.Message.Text, "/stats")))
		if e != nil {
			response = "Request failed"
		} else {
			if appConfig.statsAuth() && (fakeURL.Query().Get("username") != appConfig.StatsUsername || fakeURL.Query().Get("password") != appConfig.StatsPassword) {
				response = "Not authorized. Add username=yourusername&password=yourpassword to the query."
			} else {
				result, e := doRequest(fakeURL.Query())
				if e != nil {
					response = "Request failed"
				} else {
					rowStrings := make([]string, len(result))
					for i, row := range result {
						rowStrings[i] = (*row).First + ": " + strconv.Itoa((*row).Second)
					}
					response = strings.Join(rowStrings, "\n")
				}
			}
		}
		e = app.telegram.replyToMessage(u.Message.Chat.Id, response, u.Message.Id)
		if e != nil {
			fmt.Println("Failed to send message:", e)
		}
	}
}
