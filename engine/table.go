package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Windows holds sliding or tumbling window analysis
// where map[int] key is an integer marking the beginning of a time window in milliseconds
// and map[string] key is the count for that metric
// float64 is count or sum
type Windows map[int]map[string]float64

// Table struct is used to represent a table instance.
type Table struct {
	database    string
	file        *os.File      // a handle to the open .godbd data file
	commands    chan *Command // a channel for client to push incoming records
	thingsToMap chan JSON
	config      JSON
}

// NewTable is used to open a .godbd file.
// It returns a *Table struct
func NewTable(database, path string, config JSON) *Table {

	file, tableFileOpenErr := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)

	if tableFileOpenErr != nil {
		panic(tableFileOpenErr)
	}

	table := Table{database, file, make(chan *Command, 10), make(chan JSON, 10), config}
	//fmt.Printf("config %v\n", config)
	return &table
}

// Create writes a record to the table's .godbd file.
// It returns an error and status code
func (table *Table) Create(data JSON) (uint16, error) {

	// json encode reply into payload
	jsonString, dataMarshalErr := json.Marshal(data)
	if dataMarshalErr != nil {
		fmt.Printf("dataMarshalErr %v", dataMarshalErr)
		return 500, errors.New("ERR_ENCODE_FAILED")
	}

	newLine := []byte("\n")
	jsonLine := append(jsonString, newLine[0])
	//fmt.Printf("%s", jsonLine)
	table.file.Write(jsonLine)

	return 201, nil
}

/*
Run iterates over a table's command channel looking for actions to perform
*/
func (table *Table) executeCommands() {

	for command := range table.commands {
		switch command.Action {
		case "c":
			table.thingsToMap <- command.Data
			var err error
			reply := NewReply(command.ID)
			reply.Status, err = table.Create(command.Data)
			if err != nil {
				reply.Error = err.Error()
			}

			command.connection.replies <- reply

		}
	}
}

func (table *Table) executeMaps() {

	win := make(Windows)
	var now int
	var key string
	var intrvl int
	//var value float64
	var mp map[string]interface{}

	for thingToMap := range table.thingsToMap {
		now = int(thingToMap["time"].(float64))
		for _, m := range table.config {
			mp = m.(map[string]interface{})
			intrvl = int(mp["Interval"].(float64))
			beginningOfWindow, didMakeNewWindow := interval(win, now, intrvl)

			for _, emit := range mp["Emit"].([]interface{}) {
				e := emit.(map[string]interface{})
				keyByValue := e["Key"].(map[string]interface{})["ByValue"]
				operation := e["Value"].(map[string]interface{})["ByOperation"]

				if keyByValue != nil {
					key = keyByValue.(string)
				}

				switch operation {
				case "OP:COUNT":
					increment(win, beginningOfWindow, thingToMap[key].(string))
				case "OP:SUM":
					sum(win, beginningOfWindow, key, thingToMap[key].(float64))
				case "OP:AVG":
					//avg(win, beginningOfWindow)
				}

			}

			if didMakeNewWindow {
				fmt.Printf("%v %v\n", beginningOfWindow-intrvl, win[beginningOfWindow-intrvl])
			}

		}
	}
}

func interval(win Windows, timestamp, interval int) (int, bool) {
	didMakeNewWindow := false
	beginningOfWindow := timestamp - (timestamp % interval)

	if win[beginningOfWindow] == nil {
		win[beginningOfWindow] = make(map[string]float64)
		didMakeNewWindow = true
	}
	return beginningOfWindow, didMakeNewWindow
}

func increment(win Windows, beginningOfWindow int, field string) {
	previousValue := win[beginningOfWindow][field]
	win[beginningOfWindow][field] = previousValue + 1
}

func sum(win Windows, beginningOfWindow int, field string, value float64) {
	win[beginningOfWindow][field] = win[beginningOfWindow][field] + value
}
