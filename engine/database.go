/*
Package engine provides access to the godb core engine.
*/
package engine

import (
	"fmt"
	"os"
	"strings"
	//"errors"
	"path/filepath"
)

/*
LoadDatabases scans dataPath returning map[string]map[string]*Table
where the first map[string] key is the name of the database
and the second map[string] keys are names of tables in that database
*/
func LoadDatabases(dataPath string) map[string]map[string]*Table {
	// LoadDatabases iterates the specified path loading a db per directory
	// It returns a map of Db structs
	databases := make(map[string]map[string]*Table)

	filepath.Walk(dataPath, func(path string, f os.FileInfo, err error) error {

		if f.IsDir() && dataPath != path {
			fmt.Printf("initializing database: %s\n", f.Name())
			databases[f.Name()] = NewDatabase(f.Name(), dataPath+"/"+f.Name())
		}

		return nil
	})

	return databases

}

/*
NewDatabase traverses path assuming that all files in the folder are .godbd data files
returns a map[string]*Table where the map key is the name of the table taken from the name of the data file
*/
func NewDatabase(name, path string) map[string]*Table {

	database := make(map[string]*Table)

	// get  list of  tables

	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {

		if !f.IsDir() {
			tableName := strings.Split(f.Name(), ".")[0]

			fmt.Printf("opening table %s @ %d bytes\n", tableName, f.Size())
			tbl := NewTable(name, path)
			// TODO see if better place for this goroutine
			go tbl.Run()
			database[tableName] = tbl
		}

		return nil
	})

	return database
}
