package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

type config struct {
	Port          string `json:"port"`
	Dnt           bool   `json:"dnt"`
	DbPath        string `json:"dbPath"`
	StatsUsername string `json:"statsUsername"`
	StatsPassword string `json:"statsPassword"`
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

func (ac *config) statsAuth() bool {
	return len(ac.StatsUsername) > 0 && len(ac.StatsPassword) > 0
}
