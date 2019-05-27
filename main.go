package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	tgBot     *tgbotapi.BotAPI
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
	initTelegramBot()
}

func main() {
	go startListeningToWeb()
	go startReports()
	go startStatsTelegram()
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

func initTelegramBot() {
	if appConfig.TGBotToken == "" {
		fmt.Println("Telegram not configured.")
		return
	}
	bot, e := tgbotapi.NewBotAPI(appConfig.TGBotToken)
	if e != nil {
		fmt.Println("Failed to setup Telegram:", e)
		return
	}
	fmt.Println("Authorized Telegram bot on account", bot.Self.UserName)
	app.tgBot = bot
}

func startListeningToWeb() {
	port := appConfig.Port
	addr := ":" + port
	fmt.Printf("Listening to %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.router))
}
