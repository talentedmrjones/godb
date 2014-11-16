package engine

import (
	"os"
)

// Db struct is used to represent a database instance.
type Db struct {
	tables map[string]*Table			// collection of tables
}


// Table struct is used to represent a table instance.
type Table struct  {
	chunkSize uint32							// number of bytes each record will use on disk
	tableFileSize int64						// the file size of the table for tracking EOF/positioning for index
	tableFile *os.File						// a handle to the open .godbd data file
	primaryIndex map[string]int64	// the in-memory map of primary keys (id) -> position of chunk
	freeChunks []int64		// collection of chunks that will be reused for creates
}
