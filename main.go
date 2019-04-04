package main

import (
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"kis3.dev/kis3/helpers"
	"log"
	"math"
	"net/http"
	"strconv"
)

type kis3 struct {
	db        *Database
	router    *mux.Router
	staticBox *packr.Box
	staticFS  http.Handler
}

var (
	app = &kis3{}
)

func init() {
	e := app.setupDB()
	if e != nil {
		log.Fatal("Database setup failed:", e)
	}
	app.setupRouter()
}

func main() {
	app.startListening()
}

func (kis3 *kis3) setupDB() (e error) {
	kis3.db, e = initDatabase()
	return
}

func (kis3 *kis3) setupRouter() {
	kis3.router = mux.NewRouter()

	corsHandler := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))

	viewRouter := kis3.router.PathPrefix("/view").Subrouter()
	viewRouter.Use(corsHandler)
	viewRouter.Path("").HandlerFunc(kis3.trackView)

	kis3.router.HandleFunc("/stats", kis3.requestStats)

	staticRouter := kis3.router.PathPrefix("").Subrouter()
	staticRouter.Use(corsHandler)
	staticRouter.PathPrefix("").Handler(http.HandlerFunc(kis3.serveStaticFile))
}

func (kis3 kis3) startListening() {
	port := appConfig.port
	addr := ":" + port
	fmt.Printf("Listening to %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, kis3.router))
}

func (kis3 kis3) trackView(w http.ResponseWriter, r *http.Request) {
	url := r.Header.Get("Referer") // URL of requesting source
	ref := r.URL.Query().Get("ref")
	if r.Header.Get("DNT") == "1" && appConfig.dnt {
		fmt.Println("Not tracking because of DNT")
	} else {
		fmt.Printf("Tracking %s with referrer %s\n", url, ref)
		go kis3.db.trackView(url, ref) // run with goroutine for awesome speed!
		_, _ = fmt.Fprint(w, "true")
	}
}

func sendHelloResponse(w http.ResponseWriter) {
	_, _ = fmt.Fprint(w, "Hello from KISSS")
}

func (kis3 kis3) serveStaticFile(w http.ResponseWriter, r *http.Request) {
	if kis3.staticBox == nil || kis3.staticFS == nil {
		kis3.staticBox = packr.New("staticFiles", "./static")
		kis3.staticFS = http.FileServer(kis3.staticBox)
	}
	uPath := r.URL.Path
	if uPath != "/" && kis3.staticBox.Has(uPath) {
		kis3.staticFS.ServeHTTP(w, r)
	} else {
		sendHelloResponse(w)
	}
}

func (kis3 kis3) requestStats(w http.ResponseWriter, r *http.Request) {
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
	case "hours":
		view = HOURS
	case "days":
		view = DAYS
	case "weeks":
		view = WEEKS
	case "months":
		view = MONTHS
	}
	result, e := kis3.db.request(&ViewsRequest{
		view:   view,
		from:   queries.Get("from"),
		to:     queries.Get("to"),
		url:    queries.Get("url"),
		domain: queries.Get("domain"),
		ref:    queries.Get("ref"),
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
	var values []chart.Value
	max := float64(1)
	for _, row := range result {
		values = append(values, chart.Value{Label: row.First, Value: float64(row.Second), Style: chart.Style{
			FillColor:   drawing.ColorBlue,
			StrokeColor: drawing.ColorBlue,
		}})
		max = math.Max(max, float64(row.Second))
	}
	chartRange := &chart.ContinuousRange{
		Min: float64(0),
		Max: max,
	}
	chartWidth := len(values)*30 + 100
	barChart := chart.BarChart{
		Title:      "Stats",
		Height:     500,
		Width:      chartWidth,
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		BarWidth:   20,
		BarSpacing: 10,
		XAxis: chart.Style{
			Show:                true,
			TextRotationDegrees: 90.0,
		},
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
			Range: chartRange,
		},
		Bars: values,
	}
	w.Header().Set("Content-Type", chart.ContentTypeSVG)
	e := barChart.Render(chart.SVG, w)
	if e != nil {
		sendPlainResponse(result, w) // Fallback to plain
	}
}
