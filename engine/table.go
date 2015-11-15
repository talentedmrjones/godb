package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Table struct is used to represent a table instance.
type Table struct {
	database string
	file     *os.File      // a handle to the open .godbd data file
	commands chan *Command // a channel for client to push incoming records
}

// NewTable is used to open a .godbd file.
// It returns a *Table struct
func NewTable(database, path string) *Table {

	file, tableFileOpenErr := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)

	if tableFileOpenErr != nil {
		panic(tableFileOpenErr)
	}

	table := Table{database, file, make(chan *Command)}

	return &table
}

// Create writes a record to the table's .godbd file.
// It returns an error and status code
func (table *Table) Create(data map[string]interface{}) (uint16, error) {

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
Run loops over a table's command channel looking for actions to perform
*/
func (table *Table) Run() {

	for command := range table.commands {
		switch command.Action {
		case "c":

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
