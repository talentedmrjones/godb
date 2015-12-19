package main

import (
	"github.com/talentedmrjones/godb/engine"
	"log"
)

func main() {
	config, configError := engine.LoadConfig("./config/maps.json")
	if configError != nil {
		log.Fatalf("config error %s", configError)
	} else {
		databases := engine.LoadDatabases("./data", config)
		engine.NewServer(databases)
	}
}
