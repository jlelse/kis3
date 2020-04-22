package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

type kis3 struct {
	router    *mux.Router
	staticBox *packr.Box
	telegram  *telegram
}

var (
	app = &kis3{
		staticBox: packr.New("staticFiles", "./static"),
	}
)

func main() {
	// Init
	initConfig()
	e := initDatabase()
	if e != nil {
		log.Fatal("Database setup failed:", e)
	}
	initTelegram()
	initRouter()
	// Start
	go startListeningToWeb()
	go startReports()
	// Graceful stop
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, os.Interrupt, syscall.SIGTERM)
	sig := <-gracefulStop
	fmt.Printf("Received signal: %+v", sig)
	os.Exit(0)
}

func initRouter() {
	app.router = mux.NewRouter()
	initStatsRouter()
	initTrackingRouter()
	if app.telegram != nil {
		initTelegramRouter()
	}
}

func startListeningToWeb() {
	port := appConfig.Port
	addr := ":" + port
	fmt.Printf("Listening to %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.router))
}
