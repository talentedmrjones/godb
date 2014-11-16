package engine

import (
	"os"
)

// Db struct is used to represent a database instance.
type Db struct {
	Tables map[string]*Table
}


// Table struct is used to represent a table instance.
type Table struct  {
	TableFileSize int64
	TableFile *os.File
	//IndicesFileSize int
	PrimaryIndex map[string][2]int64
	//SecondaryIndices map[string]string
	//indicesHandle *os.File
}
