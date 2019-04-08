package main

import (
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"html/template"
	"kis3.dev/kis3/helpers"
	"log"
	"net/http"
	"strconv"
)

type kis3 struct {
	db        *Database
	router    *mux.Router
	staticBox *packr.Box
}

var (
	app = &kis3{}
)

func init() {
	e := setupDB()
	if e != nil {
		log.Fatal("Database setup failed:", e)
	}
	setupRouter()
	app.staticBox = packr.New("staticFiles", "./static")
}

func main() {
	startListening()
}

func setupDB() (e error) {
	app.db, e = initDatabase()
	return
}

func setupRouter() {
	app.router = mux.NewRouter()

	corsHandler := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))

	viewRouter := app.router.PathPrefix("/view").Subrouter()
	viewRouter.Use(corsHandler)
	viewRouter.Path("").HandlerFunc(trackView)

	app.router.HandleFunc("/stats", requestStats)

	staticRouter := app.router.PathPrefix("").Subrouter()
	staticRouter.Use(corsHandler)
	staticRouter.HandleFunc("/kis3.js", serveTrackingScript)
	staticRouter.PathPrefix("").Handler(http.HandlerFunc(sendHelloResponse))
}

func startListening() {
	port := appConfig.port
	addr := ":" + port
	fmt.Printf("Listening to %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.router))
}

func trackView(w http.ResponseWriter, r *http.Request) {
	url := r.Header.Get("Referer") // URL of requesting source
	ref := r.URL.Query().Get("ref")
	ua := r.Header.Get("User-Agent")
	if !(r.Header.Get("DNT") == "1" && appConfig.dnt) {
		go app.db.trackView(url, ref, ua) // run with goroutine for awesome speed!
		_, _ = fmt.Fprint(w, "true")
	}
}

func sendHelloResponse(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "Hello from KISSS")
}

func serveTrackingScript(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "application/javascript")
	trackingScriptBytes, _ := app.staticBox.Find("kis3.js")
	_, _ = w.Write(trackingScriptBytes)
}

func requestStats(w http.ResponseWriter, r *http.Request) {
	// Require authentication
	if appConfig.statsAuth {
		if !helpers.CheckAuth(w, r, appConfig.statsUsername, appConfig.statsPassword) {
			return
		}
	}
	// Do request
	queries := r.URL.Query()
	view := PAGES
	switch queries.Get("view") {
	case "pages":
		view = PAGES
	case "referrers":
		view = REFERRERS
	case "useragents":
		view = USERAGENTS
	case "hours":
		view = HOURS
	case "days":
		view = DAYS
	case "weeks":
		view = WEEKS
	case "months":
		view = MONTHS
	}
	result, e := app.db.request(&ViewsRequest{
		view: view,
		from: queries.Get("from"),
		to:   queries.Get("to"),
		url:  queries.Get("url"),
		ref:  queries.Get("ref"),
		ua:   queries.Get("ua"),
	})
	if e != nil {
		fmt.Println("Database request failed:", e)
		w.WriteHeader(500)
	} else if result != nil {
		switch queries.Get("format") {
		case "json":
			sendJsonResponse(result, w)
			return
		case "chart":
			sendChartResponse(result, w)
			return
		default: // "plain"
			sendPlainResponse(result, w)
			return
		}
	}
}

func sendPlainResponse(result []*RequestResultRow, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "text/plain")
	for _, row := range result {
		_, _ = fmt.Fprintln(w, (*row).First+": "+strconv.Itoa((*row).Second))
	}
}

func sendJsonResponse(result []*RequestResultRow, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
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
