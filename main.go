package main

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type kis3 struct {
	router    *mux.Router
	staticBox *packr.Box
}

var (
	app = &kis3{
		staticBox: packr.New("staticFiles", "./static"),
	}
)

func init() {
	initConfig()
	e := initDatabase()
	if e != nil {
		log.Fatal("Database setup failed:", e)
	}
	initRouter()
}

func main() {
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
}

func startListeningToWeb() {
	port := appConfig.Port
	addr := ":" + port
	fmt.Printf("Listening to %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.router))
}
