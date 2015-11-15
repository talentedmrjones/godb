package engine

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// Connection is a single tcp connection
type Connection struct {
	socket  net.Conn   // holds the network socket
	replies chan Reply // handles replies sent to client
}

/*
NewConnection creates a Connection
*/
func NewConnection(conn net.Conn) *Connection {
	// initialize an instance of the Connection struct
	return &Connection{conn, make(chan Reply)}

}

// Receive continuously looks for data from the socket and relays that to a table's command channel
func (connection *Connection) receive(databases map[string]map[string]*Table) {
	// create a buffered reader
	buffer := bufio.NewReader(connection.socket)

	// loop forever
	for {

		// prepare a buffer to read the data
		var payloadBytes []byte
		var brackets = 0
		// read until we get an even number of brackets
		for {
			// read data
			nextByte, bufferReadErr := buffer.ReadByte()
			if bufferReadErr != nil {
				// this will mean the connection is closed
				return
			}
			// add the byte to the array
			payloadBytes = append(payloadBytes, nextByte)
			// convert to string
			character := string(nextByte)

			// count brackets
			if character == "{" {
				brackets++
			} else if character == "}" {
				brackets--
			}

			if brackets == 0 && len(payloadBytes) > 1 {
				break
			}
		}

		if len(payloadBytes) > 1 {
			//fmt.Printf("%s", payloadBytes)

			// prepare a command struct to hold incoming JSON
			command := NewCommand(connection)

			payloadBytesUnmarshalErr := json.Unmarshal(payloadBytes, command)
			if payloadBytesUnmarshalErr != nil {
				log.Fatal("payloadBytesUnmarshalErr decode:", payloadBytesUnmarshalErr)
				// TODO: report error to connection
			}

			// deliver data to databases table
			databases[command.Db][command.Table].commands <- command

		}
	}
}

// send is run in its own goroutine. It continuously loops over replies channel handling replies for that connection
func (connection *Connection) send() {
	for reply := range connection.replies {

		//fmt.Printf("%#v\n", reply)
		// json encode reply into payload
		payloadBytes, replyMarshalErr := json.Marshal(reply)
		if replyMarshalErr != nil {
			fmt.Printf("commandMarshalError %v", replyMarshalErr)
		}

		dataSize := make([]byte, 4)
		binary.BigEndian.PutUint32(dataSize, uint32(len(payloadBytes)))

		connection.socket.Write(dataSize)
		connection.socket.Write(payloadBytes)
	}
}
