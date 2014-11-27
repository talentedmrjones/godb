package engine

import (
	"net"
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
)



// Client is a single tcp connection
type Client struct {
	conn		net.Conn 		// holds the network socket
	//ch 			chan string			// receives replies/errors from table
}


// ReadLinesInto continuously looks for data from the connection and relays that to the message channel
func (c *Client) Receive(databases map[string]map[string]*Table) {
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
		// first byte is one of: c,r,u,d
		command := string(header[0:1])
		// last 4 bytes are 32 bit unsigned int for size of data to follow
		dataSize := binary.BigEndian.Uint32(header[1:5])

		// prepare a buffer to read the data
		payloadBytes := make([]byte, dataSize)
		// read data
		numDataBytes, err := buf.Read(payloadBytes)
		if uint32(numDataBytes)<dataSize || err != nil {
			// report error here
			break
		}


		// Create a decoder and receive a value.
		payloadDecoderBuffer := bytes.NewBuffer(payloadBytes)
		payloadDecoder := gob.NewDecoder(payloadDecoderBuffer)
		payload := make(map[string][]byte)
		err = payloadDecoder.Decode(&payload)
		if err != nil {
			log.Fatal("decode:", err)
		}


		database := string(payload["db"])
		table := string(payload["tbl"])

		fmt.Printf("command %s on %s.%s with data:%v", command, database, table, payload["data"])
		// deliver data to databases[database][table]
		//databases[db][tbl].CommandChan<- &Command{command, bytedata, c}

		// push the line received above onto the channel
		//msgchan <- fmt.Sprintf("%s", line)
	}
}
