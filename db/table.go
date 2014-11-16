package engine

import (
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

// Table.Read reads a record from the table's .godbd file
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
  buf := make([]byte, primaryIndex[1])

  // read starting at the offset stored at index[0]
  _, readErr := tbl.TableFile.ReadAt(buf, primaryIndex[0])
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

// Table.Create writes a record to the table's .godbd file.
// It returns an error|nil or on success a map[string]string (the data that it was original given)
func (tbl *Table) Create (data map[string]string) (error, map[string]string) {

  if _, ok := data["id"]; !ok {
    err := errors.New("ID_MISSING")
    return err, nil
  }

  if _, piok := tbl.PrimaryIndex[data["id"]]; piok {
    // id exists in index, therefore not new
    return errors.New("ID_NON_UNIQUE"), nil
  }

  b := new(bytes.Buffer)
  e := gob.NewEncoder(b)

  // encode data
  err := e.Encode(data)
  if err != nil {
      panic(err)
  }

  // initialize primary index as array
  // record position,length in index
  tbl.PrimaryIndex[data["id"]] = [2]int64{tbl.TableFileSize, int64(b.Len())}

  tbl.TableFile.WriteAt(b.Bytes(), tbl.TableFileSize)

  //fmt.Printf("%#v", tbl.PrimaryIndex[data["id"]])

  setTableFileSize(tbl, b.Len())

  return nil, nil
}
