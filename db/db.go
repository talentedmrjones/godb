// Package db provides acces to the godb core engine.
package engine

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"
)

// OpenTable is used to open a .godbd file.
// It returns a *Table struct
func OpenTable (path string, tableFileSize int64) (*Table) {

	tableFile, tableFileOpenErr := os.OpenFile(path, os.O_RDWR, 0)

	if tableFileOpenErr != nil {
		panic(tableFileOpenErr)
	}

	table := Table{tableFileSize, tableFile, make(map[string][2]int64)}

	return &table
}

// NewDatabase creates a Db which contains a map of *Table
// It returns *Db struct
func NewDatabase () (*Db) {

	database := Db{make(map[string]*Table)}

	// get  list of  tables

	filepath.Walk("./data", func (path string, f os.FileInfo, err error) error {

		if (!f.IsDir()) {
			tableName := strings.Split(f.Name(),".")

			fmt.Printf("opening %s with %d bytes\n", tableName[0], f.Size())
			tbl := OpenTable(path, f.Size())

			database.Tables[tableName[0]] = tbl
		}

		return nil
	})

	return &database
}
