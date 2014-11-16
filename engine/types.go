package engine

import (
	"os"
)

// Db struct is used to represent a database instance.
type Db struct {
	tables map[string]*Table
}


// Table struct is used to represent a table instance.
type Table struct  {
	chunkSize uint32
	tableFileSize int64
	tableFile *os.File
	//IndicesFileSize int
	primaryIndex map[string]int64
	//SecondaryIndices map[string]string
	//indicesHandle *os.File
	freeChunks map[int64]uint8
}
