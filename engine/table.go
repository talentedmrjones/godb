package engine

import (
  "os"
  "encoding/gob"
  "bytes"
  "errors"
  "fmt"
)

// Table struct is used to represent a table instance.
type Table struct  {
  chunkSize 			uint32							// number of bytes each record will use on disk
  tableFileSize 	int64								// the file size of the table for tracking EOF/positioning for index
  tableFile 			*os.File						// a handle to the open .godbd data file
  // TODO look into map[[]byte]int64
  primaryIndex 		map[string]int64		// the in-memory map of primary keys (id) -> position of chunk
  freeChunks 			[]int64							// collection of chunks that will be reused for creates
  Chan 		         chan *Command				// a channel for client to push incoming records
}

// NewTable is used to open a .godbd file.
// It returns a *Table struct
func NewTable (chunkSize uint32, path string, tableFileSize int64) *Table {

  tableFile, tableFileOpenErr := os.OpenFile(path, os.O_RDWR, 0)

  if tableFileOpenErr != nil {
    panic(tableFileOpenErr)
  }

  table := Table{chunkSize, tableFileSize, tableFile, make(map[string]int64), make([]int64, 0, 5), make(chan *Command)}

  return &table
}

// grow expands the file size of the table as records are appends
// The table's file size is used to demarcate EOF which is used for indexing the position of records.
func (tbl *Table) grow (newSize int) {
  tbl.tableFileSize += int64(newSize)
}

// getFreeChunk returns the position of a chunk that was previously marked as deleted so that Create can use disk most efficiently
func (tbl *Table) getFreeChunk () int64 {

  var chunkAddress int64 = -1

  if len(tbl.freeChunks) > 0 {
    chunkAddress = tbl.freeChunks[0]
    tbl.freeChunks = tbl.freeChunks[1:]
  }

  return chunkAddress
}

// addFreeChunk adds a position to the freeChunks slice
func (tbl *Table) addFreeChunk (position int64) {
  tbl.freeChunks = append(tbl.freeChunks, position)
}


// Create writes a record to the table's .godbd file.
// It returns an error|nil or on success a map[string][]byte (the data that it was original given)
func (tbl *Table) Create (data map[string][]byte) (error, map[string][]byte) {

  // ensure id field exists
  id, dataIdExists := data["id"]
  if !dataIdExists {
    err := errors.New("ID_MISSING")
    return err, nil
  }

  if _, primaryIndexExists := tbl.primaryIndex[string(id)]; primaryIndexExists {
    // id exists in index, therefore not new
    return errors.New("ID_NON_UNIQUE"), nil
  }

  // encode data
  dataBuffer := new(bytes.Buffer)
  e := gob.NewEncoder(dataBuffer)
  err := e.Encode(data)
  if err != nil {
      return errors.New("ENCODE_FAILED"), nil
  }

  bufferLength := dataBuffer.Len()

  // ensure that the data is not larger than table chunk size
  // TODO: test

  if uint32(bufferLength)>tbl.chunkSize {
    bufferLengthErr :=  errors.New("TOO_LARGE")
    return bufferLengthErr, nil
  }

  //  pad the buffer out to the chunkSize
  dataBuffer.Write(bytes.Repeat([]byte{0}, int(tbl.chunkSize)-bufferLength))

  // determine position by looking to freeChunks
  var position int64
  position = tbl.getFreeChunk()

  if position < 0 {
    position = tbl.tableFileSize
    tbl.grow(dataBuffer.Len())
  }

  // initialize primary index as array
  // record position of record in index
  tbl.primaryIndex[string(data["id"])] = position

  tbl.tableFile.WriteAt(dataBuffer.Bytes(), position)

  return nil, data
}


// Read reads a record from the table's .godbd file
// It returns []byte
func (tbl *Table) Read (data map[string][]byte) (error, map[string][]byte) {

  // data must include an id
  if _,idok := data["id"]; !idok {
    err := errors.New("ID_MISSING")
    return err,nil
  }

  // find the offset and byte length from index
  id := string(data["id"])
  if _, indexOk := tbl.primaryIndex[id]; !indexOk {
    // wasn't found
    err := errors.New("NOT_FOUND")
    return err,nil
  }

  primaryIndex := tbl.primaryIndex[id]

  // make a []byte using the length stored at index[1]
  buf := make([]byte, tbl.chunkSize)

  // read starting at the offset stored at index[0]
  _, readErr := tbl.tableFile.ReadAt(buf, primaryIndex)
  if readErr != nil {
    panic(readErr)
  }

  // decode the buffer
  bufReader := bytes.NewReader(buf)
  d := gob.NewDecoder(bufReader)
  var decodedMap map[string][]byte

  // decode data
  decodeErr := d.Decode(&decodedMap)
  if decodeErr != nil {
      panic(decodeErr)
  }

  return nil, decodedMap
}

// Updates replaces a record to the table's .godbd file.
// It returns an error or on success a map[string][]byte (the data that it was original given)
func (tbl *Table) Update (data map[string][]byte) (error, map[string][]byte) {

  // ensure id field exists
  if _, dataIdExists := data["id"]; !dataIdExists {
    err := errors.New("ID_MISSING")
    return err, nil
  }

  // ensure record index contains the record
  if _, indexOk := tbl.primaryIndex[string(data["id"])]; !indexOk {
    // wasn't found
    err := errors.New("NOT_FOUND")
    return err, nil
  }

  // encode data
  b := new(bytes.Buffer)
  e := gob.NewEncoder(b)
  err := e.Encode(data)
  if err != nil {
      return errors.New("ENCODE_FAILED"), nil
  }

  bufferLength := b.Len()

  // ensure that the data is not larger than table chunk size
  // TODO: test

  if uint32(bufferLength)>tbl.chunkSize {
    bufferLengthErr :=  errors.New("TOO_LARGE")
    return bufferLengthErr, nil
  }

  //  pad the buffer out to the chunkSize
  b.Write(bytes.Repeat([]byte{0}, int(tbl.chunkSize)-bufferLength))

  tbl.tableFile.WriteAt(b.Bytes(), tbl.primaryIndex[string(data["id"])])

  return nil, data
}


// Delete removes the record from the index and marks the chunk for reuse
func (tbl *Table) Delete (data map[string][]byte) error {

  // ensure id field exists
  if _, dataIdExists := data["id"]; !dataIdExists {
    err := errors.New("ID_MISSING")
    return err
  }

  // ensure record index contains the record
  if _, indexOk := tbl.primaryIndex[string(data["id"])]; !indexOk {
    // wasn't found
    return nil
  }

  // mark chunk available for reuse
  tbl.addFreeChunk(tbl.primaryIndex[string(data["id"])])

  // delete record from the index
  delete(tbl.primaryIndex, string(data["id"]))

  return nil

}

func (table *Table) Run () {

  for command := range table.Chan {
    switch command.Action {
      case "c":
        // TODO send data back to client
          createErr, _ := table.Create(command.Data)
          if createErr != nil {
            // TODO handle error
            fmt.Printf("Create error %s\n", createErr)
          } else {
            fmt.Printf("Record created\n")
          }
    }
	}
}
