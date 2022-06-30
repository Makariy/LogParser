package main

import (
	"configparser"
	"fmt"
	"log"
	"logparser"
	"os"
	"units"
)

func getConfig(filename string) configparser.Config {
	configFile, err := os.Open("config.conf")
	if err != nil {
		log.Fatal("Cannot open config file")
	}
	config := configparser.Config{}
	if config.ParseConfig(configFile) != nil {
		log.Fatalf("An error occurred during parsing the config file: %s", err)
	}
	return config
}

func getUnits(config configparser.Config) []units.Unit {
	logFile, err := os.Open(config.Input.Filename)
	if err != nil {
		log.Fatal("Cannot open file")
	}
	items, err := logparser.ParseFile(logFile)
	if err != nil {
		log.Fatal("Cannot parse log file")
	}
	return items
}

func main() {
	config := getConfig("config.conf")
	items := getUnits(config)

	statistics := units.UnitStats{Units: items, Config: config}
	statistics.CalculateStats(config)

	fmt.Println(statistics.Stats["88.216.205.14"].AvgRequestsPerSecond)
	fmt.Println(statistics.Stats["88.216.205.14"].MostRequestedPaths)
	fmt.Println(statistics.MostAnnoyingIPs[:2])
}
