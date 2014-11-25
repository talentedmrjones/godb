package server

import (
	"net"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	//"io"
)



// Client is a single tcp connection
type Client struct {
	conn		net.Conn 		// holds the network socket
	ch 			chan string			// receives replies from table and errors from Client
}


// ReadLinesInto continuously looks for data from the connection and relays that to the message channel
func (c Client) Receive() {
	// create a buffered reader
	buf := bufio.NewReader(c.conn)

	// loop forever
	for {

		// read the first 5 bytes
		header := make([]byte,5)
		numHeaderBytes, err := buf.Read(header)
		if numHeaderBytes<5 || err != nil {
			// report error back to client
			break
		}

		// we should now have 5 bytes
		cmd := string(header[0:1]) // first byte is one of: c,r,u,d
		dataSize := binary.BigEndian.Uint32(header[1:5]) // last 4 bytes represents the datasize

		dataBytes := make([]byte, dataSize) // prepare a buffer to read the data
		numDataBytes,err := buf.Read(dataBytes)
		if uint32(numDataBytes)<dataSize || err != nil {
			// report error here
			break
		}

		// dataBytes will be JSON so unmarshal it
		var data map[string]interface{}
		json.Unmarshal(dataBytes, &data)
		fmt.Printf("%s on %s.%s %v\n", cmd, data["database"], data["table"], data["data"])
		// deliver data to databases[database][table]

		// push the line received above onto the channel
		//msgchan <- fmt.Sprintf("%s", line)
	}
}
