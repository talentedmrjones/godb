package engine

import (
	"os"
)

// Table struct is used to represent a table instance.
type Table struct  {
	chunkSize 			uint32							// number of bytes each record will use on disk
	tableFileSize 	int64								// the file size of the table for tracking EOF/positioning for index
	tableFile 			*os.File						// a handle to the open .godbd data file
	primaryIndex 		map[string]int64		// the in-memory map of primary keys (id) -> position of chunk
	freeChunks 			[]int64							// collection of chunks that will be reused for creates
	CommandChan 		chan *Command				// a channel for client to push commands
}

type Command struct {
	Action string 						// create,read,update,delete
	Data		map[string][]byte // the data send by client, being acted upon by the command
	Client  *Client 					// reference to the client so the Table can push to its reply channel
}

// type Record interface {
// 	Bytes() []byte // return raw data
// }
