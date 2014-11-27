package main

import (
  //"fmt"
  //"time"
  "github.com/talentedmrjones/godb/engine"
)

//var database *engine.Db
//var users *engine.Table

func main () {
  databases := engine.LoadDatabases("./data")
  engine.RunServer(databases)
  // initialize a database engine
  // database = engine.NewDatabase()
  // _, users = database.GetTable("users")
  // create()
  // read()
  // update()
  // read()
  // delete()
  // read()
  // create()
  // read()
}

// func create() {
//   start := time.Now()
//   data := map[string]string{"id":"123","name":"Richard"}
//   err,data := users.Create(data)
//
//   if err!=nil {
//     println(string(err.Error()))
//   }
//   elapsed := time.Since(start)
//   fmt.Printf("create took %v sec\n", elapsed.Seconds())
// }
//
// func read() {
//   start := time.Now()
//   query := map[string]string{"id":"123"}
//
//   err,user := users.Read(query)
//
//   if err!=nil {
//     println(string(err.Error()))
//   }
//   elapsed := time.Since(start)
//   fmt.Printf("%s %s in %f sec\n", user["id"], user["name"], elapsed.Seconds())
//
// }
//
// func update() {
//   start := time.Now()
//   data := map[string]string{"id":"123","name":"Mike"}
//   err,data := users.Update(data)
//
//   if err!=nil {
//     println(string(err.Error()))
//   }
//   elapsed := time.Since(start)
//   fmt.Printf("update took %v sec\n", elapsed.Seconds())
// }
//
// func delete() {
//   start := time.Now()
//   data := map[string]string{"id":"123"}
//   err := users.Delete(data)
//   if err!=nil {
//     println(string(err.Error()))
//   }
//
//   elapsed := time.Since(start)
//   fmt.Printf("delete took %v sec\n", elapsed.Seconds())
//
// }
