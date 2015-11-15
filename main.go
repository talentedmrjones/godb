package main

import (
	"github.com/talentedmrjones/godb/engine"
)

func main() {
	databases := engine.LoadDatabases("./data")
	engine.NewServer(databases)
}
