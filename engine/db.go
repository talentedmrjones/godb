// Package db provides acces to the godb core engine.
package engine

import (
	"os"
	"fmt"
	"strings"
	"errors"
	"path/filepath"
)


// NewDatabase creates a Db which contains a map of *Table
// It returns *Db struct
func NewDatabase () (*Db) {

	database := Db{make(map[string]*Table)}

	// get  list of  tables

	filepath.Walk("./data", func (path string, f os.FileInfo, err error) error {

		if (!f.IsDir()) {
			tableName := strings.Split(f.Name(),".")

			fmt.Printf("opening %s with %d bytes and %d chunk size\n", tableName[0], f.Size(), 4096)
			tbl := NewTable(4096, path, f.Size())

			database.AddTable(tableName[0], tbl)
		}

		return nil
	})

	return &database
}

// GetTable returns a *Table from the Db.tables map
func (db *Db) GetTable (tableName string) (error, *Table) {
	table, tableExists := db.tables[tableName]
	if !tableExists {
		return errors.New("TABLE_UNKNOWN"),nil
	}
	return nil, table
}

// AddTable adds a *Table to the Db.tables map
func (db *Db) AddTable (tableName string, table *Table) (error) {
	db.tables[tableName] = table
	return nil
}
