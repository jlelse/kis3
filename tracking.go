package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"net/http"
)

func initTrackingRouter() {
	corsHandler := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))
	viewRouter := app.router.Path("/view").Subrouter()
	viewRouter.Use(corsHandler)
	viewRouter.Path("").HandlerFunc(TrackingHandler)
	scriptRouter := app.router.Path("/kis3.js").Subrouter()
	scriptRouter.Use(corsHandler)
	scriptRouter.HandleFunc("", TrackingScriptHandler)
}

func TrackingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
	url := r.URL.Query().Get("url")
	ref := r.URL.Query().Get("ref")
	ua := r.Header.Get("User-Agent")
	if !(r.Header.Get("DNT") == "1" && appConfig.Dnt) {
		go trackView(url, ref, ua) // run with goroutine for awesome speed!
		_, _ = fmt.Fprint(w, "true")
	}
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
