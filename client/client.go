package main

import (
	//"fmt"
	"bufio"
	"net"
	"os"
	//"io"
	"encoding/binary"
)

func main () {
		conn, err := net.Dial("tcp", "127.0.0.1:6000")

		if err!=nil {
			// handle error here
		}

		defer conn.Close()

		for {
			console := bufio.NewReader(os.Stdin)
			data, err := console.ReadBytes('\n')
			if err!= nil {
				// handle error
			}
			//data := []byte("{\"database\":\"tmj\",\"table\":\"users\",\"data\":{\"id\":123}}")
			dataSize := make([]byte,4)
			binary.BigEndian.PutUint32(dataSize, uint32(len(data)))
			header := []byte{'c'}
			header = append(header, dataSize...)

			conn.Write(header)
			conn.Write(data)

		}

}
