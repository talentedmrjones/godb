// Package db provides acces to the godb core engine.
package engine

import (
	"os"
	"fmt"
	"strings"
	//"errors"
	"path/filepath"
)

// NewDatabase creates a Db which contains a map of *Table
// It returns *Db struct
func LoadDatabases (dataPath string) (map[string]map[string]*Table) {
	// LoadDatabases iterates the specified path loading a db per directory
	// It returns a map of Db structs
	databases := make(map[string]map[string]*Table)

	filepath.Walk(dataPath, func (path string, f os.FileInfo, err error) error {

		if (f.IsDir() && dataPath!=path) {
			fmt.Printf("initializing database: %s\n", f.Name())
			databases[f.Name()] = NewDatabase(dataPath+"/"+f.Name())
		}

		return nil
	})

	return databases

}


func NewDatabase (path string) (map[string]*Table) {

	database := make(map[string]*Table)

	// get  list of  tables

	filepath.Walk(path, func (path string, f os.FileInfo, err error) error {

		if (!f.IsDir()) {
			tableName := strings.Split(f.Name(),".")

			fmt.Printf("adding table %s %d/%d\n", tableName[0], 4096, f.Size())
			tbl := NewTable(4096, path, f.Size())
			// TODO see if better place for this goroutine
			go tbl.Run()
			database[tableName[0]]=tbl
		}

		return nil
	})

	return database
}
