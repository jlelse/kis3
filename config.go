package main

import (
	"os"
	"strconv"
)

type config struct {
	port          string
	dnt           bool
	dbPath        string
	statsAuth     bool
	statsUsername string
	statsPassword string
}

var (
	appConfig = &config{}
)

func init() {
	appConfig.port = port()
	appConfig.dnt = dnt()
	appConfig.dbPath = dbPath()
	appConfig.statsUsername = statsUsername()
	appConfig.statsPassword = statsPassword()
	appConfig.statsAuth = statsAuth(appConfig)
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) != 0 {
		return port
	} else {
		return "8080"
	}
}

func dnt() bool {
	dnt := os.Getenv("DNT")
	dntBool, e := strconv.ParseBool(dnt)
	if e != nil {
		dntBool = true
	}
	return dntBool
}

func dbPath() (dbPath string) {
	dbPath = os.Getenv("DB_PATH")
	if len(dbPath) == 0 {
		dbPath = "data/kis3.db"
	}
	return
}

func statsUsername() (username string) {
	username = os.Getenv("STATS_USERNAME")
	return
}

func statsPassword() (password string) {
	password = os.Getenv("STATS_PASSWORD")
	return
}

func statsAuth(ac *config) bool {
	return len(ac.statsUsername) > 0 && len(ac.statsPassword) > 0
}
