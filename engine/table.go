package engine

import (
  "os"
  "encoding/gob"
  "bytes"
  "errors"
  //"fmt"
)

/* NON-EXPORTED METHODS */

// setTableFileSize is used internally to update the file size of the table.
// The table's file size is used to demarcate EOF which is used for indexing.
func setTableFileSize (tbl *Table, newSize int) {
  tbl.TableFileSize += int64(newSize)
}

/* EXPORTED METHODS */

// NewTable is used to open a .godbd file.
// It returns a *Table struct
func NewTable (chunkSize uint32, path string, tableFileSize int64) (*Table) {

  tableFile, tableFileOpenErr := os.OpenFile(path, os.O_RDWR, 0)

  if tableFileOpenErr != nil {
    panic(tableFileOpenErr)
  }

  table := Table{chunkSize, tableFileSize, tableFile, make(map[string]int64), make(map[int64]uint8)}

  return &table
}


// Create writes a record to the table's .godbd file.
// It returns an error|nil or on success a map[string]string (the data that it was original given)
func (tbl *Table) Create (data map[string]string) (error, map[string]string) {

  // ensure id field exists
  if _, dataIdExists := data["id"]; !dataIdExists {
    err := errors.New("ID_MISSING")
    return err, nil
  }

  if _, primaryIndexExists := tbl.PrimaryIndex[data["id"]]; primaryIndexExists {
    // id exists in index, therefore not new
    return errors.New("ID_NON_UNIQUE"), nil
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

  if uint32(bufferLength)>tbl.ChunkSize {
    bufferLengthErr :=  errors.New("TOO_LARGE")
    return bufferLengthErr, nil
  }

  //  pad the buffer out to the chunkSize
  b.Write(bytes.Repeat([]byte{0}, int(tbl.ChunkSize)-bufferLength))

  // initialize primary index as array
  // record position of record in index
  tbl.PrimaryIndex[data["id"]] = tbl.TableFileSize

  tbl.TableFile.WriteAt(b.Bytes(), tbl.TableFileSize)

  setTableFileSize(tbl, b.Len())

  return nil, data
}


// Read reads a record from the table's .godbd file
// It returns []byte
func (tbl *Table) Read (data map[string]string) (error, map[string]string) {

  // data must include an id
  if _,idok := data["id"]; !idok {
    err := errors.New("ID_MISSING")
    return err,nil
  }

  // find the offset and byte length from index
  id := data["id"]
  if _, indexOk := tbl.PrimaryIndex[id]; !indexOk {
    // wasn't found
    err := errors.New("NOT_FOUND")
    return err,nil
  }

  primaryIndex := tbl.PrimaryIndex[id]

  // make a []byte using the length stored at index[1]
  buf := make([]byte, tbl.ChunkSize)

  // read starting at the offset stored at index[0]
  _, readErr := tbl.TableFile.ReadAt(buf, primaryIndex)
  if readErr != nil {
    panic(readErr)
  }

  // decode the buffer
  bufReader := bytes.NewReader(buf)
  d := gob.NewDecoder(bufReader)
  var decodedMap map[string]string

  // decode data
  decodeErr := d.Decode(&decodedMap)
  if decodeErr != nil {
      panic(decodeErr)
  }

  return nil, decodedMap
}

// Updates replaces a record to the table's .godbd file.
// It returns an error or on success a map[string]string (the data that it was original given)
func (tbl *Table) Update (data map[string]string) (error, map[string]string) {

  // ensure id field exists
  if _, dataIdExists := data["id"]; !dataIdExists {
    err := errors.New("ID_MISSING")
    return err, nil
  }

  // ensure record index contains the record
  if _, indexOk := tbl.PrimaryIndex[data["id"]]; !indexOk {
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

  if uint32(bufferLength)>tbl.ChunkSize {
    bufferLengthErr :=  errors.New("TOO_LARGE")
    return bufferLengthErr, nil
  }

  //  pad the buffer out to the chunkSize
  b.Write(bytes.Repeat([]byte{0}, int(tbl.ChunkSize)-bufferLength))

  tbl.TableFile.WriteAt(b.Bytes(), tbl.PrimaryIndex[data["id"]])

  setTableFileSize(tbl, b.Len())

  return nil, data
}


// Delete removes the record from the index and marks the chunk for reuse
func (tbl *Table) Delete (data map[string]string) (error) {

  // ensure id field exists
  if _, dataIdExists := data["id"]; !dataIdExists {
    err := errors.New("ID_MISSING")
    return err
  }

  // ensure record index contains the record
  if _, indexOk := tbl.PrimaryIndex[data["id"]]; !indexOk {
    // wasn't found
    return nil
  }

  // mark chunk available for reuse
  tbl.FreeChunks[tbl.PrimaryIndex[data["id"]]]=1

  // delete record from the index
  delete(tbl.PrimaryIndex,data["id"])

  return nil

}
