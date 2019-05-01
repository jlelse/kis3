package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strconv"
)

type config struct {
	Port          string   `json:"port"`
	Dnt           bool     `json:"dnt"`
	DbPath        string   `json:"dbPath"`
	StatsUsername string   `json:"statsUsername"`
	StatsPassword string   `json:"statsPassword"`
	Reports       []report `json:"reports"`
}

type report struct {
	Name         string `json:"name"`
	Time         string `json:"time"`
	To           string `json:"to"`
	SmtpUser     string `json:"smtpUser"`
	SmtpPassword string `json:"smtpPassword"`
	SmtpHost     string `json:"smtpHost"`
	From         string `json:"from"`
	Query        string `json:"query"`
}

var (
	appConfig = &config{
		Port:          "8080",
		Dnt:           true,
		DbPath:        "data/kis3.db",
		StatsUsername: "",
		StatsPassword: "",
	}
)

func init() {
	parseConfigFile(appConfig)
	// Replace values that are set via environment vars (to make it compatible with old method)
	overwriteEnvVarValues(appConfig)
}

func parseConfigFile(appConfig *config) {
	configFile := flag.String("c", "config.json", "Config file")
	flag.Parse()
	configJson, e := ioutil.ReadFile(*configFile)
	if e != nil {
		return
	}
	e = json.Unmarshal([]byte(configJson), appConfig)
	if e != nil {
		return
	}
	return
}

func overwriteEnvVarValues(appConfig *config) {
	port, set := os.LookupEnv("PORT")
	if set {
		appConfig.Port = port
	}
	dntString, set := os.LookupEnv("DNT")
	dntBool, e := strconv.ParseBool(dntString)
	if set && e == nil {
		appConfig.Dnt = dntBool
	}
	dbPath, set := os.LookupEnv("DB_PATH")
	if set {
		appConfig.DbPath = dbPath
	}
	username, set := os.LookupEnv("STATS_USERNAME")
	if set {
		appConfig.StatsUsername = username
	}
	password, set := os.LookupEnv("STATS_PASSWORD")
	if set {
		appConfig.StatsPassword = password
	}
}

func (ac *config) statsAuth() bool {
	return len(ac.StatsUsername) > 0 && len(ac.StatsPassword) > 0
}
