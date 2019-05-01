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
	"strings"
)

type kis3 struct {
	db        *Database
	router    *mux.Router
	staticBox *packr.Box
}

var (
	app = &kis3{
		staticBox: packr.New("staticFiles", "./static"),
	}
)

func init() {
	e := setupDB()
	if e != nil {
		log.Fatal("Database setup failed:", e)
	}
	setupRouter()
	setupReports()
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
	viewRouter.Path("").HandlerFunc(TrackingHandler)

	app.router.HandleFunc("/stats", StatsHandler)

	staticRouter := app.router.PathPrefix("").Subrouter()
	staticRouter.Use(corsHandler)
	staticRouter.HandleFunc("/kis3.js", TrackingScriptHandler)
	staticRouter.PathPrefix("").Handler(http.HandlerFunc(HelloResponseHandler))
}

func startListening() {
	port := appConfig.Port
	addr := ":" + port
	fmt.Printf("Listening to %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.router))
}

func TrackingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
	url := r.URL.Query().Get("url")
	ref := r.URL.Query().Get("ref")
	ua := r.Header.Get("User-Agent")
	if !(r.Header.Get("DNT") == "1" && appConfig.Dnt) {
		go app.db.trackView(url, ref, ua) // run with goroutine for awesome speed!
		_, _ = fmt.Fprint(w, "true")
	}
}

func HelloResponseHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "Hello from KISSS")
}

func TrackingScriptHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "public, max-age=432000") // 5 days
	filename := "kis3.js"
	file, err := app.staticBox.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return
	}
	http.ServeContent(w, r, filename, stat.ModTime(), file)
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
	}
	result, e := app.db.request(&ViewsRequest{
		view:     view,
		from:     queries.Get("from"),
		to:       queries.Get("to"),
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
