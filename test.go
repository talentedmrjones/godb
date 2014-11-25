package main

import (
	//"encoding/binary"
	"fmt"
	)

func main () {
	data := make(map[string]interface{})
	data["database"]="tmj"
	data["table"]="users"
	data["data"]= make(map[string]interface{})
	data["data"]["user"]="Richard"
	fmt.Printf("%v", data)
}
