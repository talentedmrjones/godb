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
	conn			net.Conn 		// holds the network socket
	replies 	chan string			// receives replies/errors from table
}

// ReadLinesInto continuously looks for data from the connection and relays that to the message channel
func (c *Client) Receive(databases map[string]map[string]*Table) {
	// create a buffered reader
	buf := bufio.NewReader(c.conn)

	// loop forever
	for {
		// read the first 4 bytes
		dataSizeBytes := make([]byte,4)
		_, err := buf.Read(dataSizeBytes)
		if err != nil {
			// report error back to client
			break
		}

		// convert those 4 bytes to 32 bit unsigned int for size of data to follow
		dataSize := binary.BigEndian.Uint32(dataSizeBytes)

		// prepare a buffer to read the data
		payloadBytes := make([]byte, dataSize)
		// read data
		numDataBytes, err := buf.Read(payloadBytes)
		if uint32(numDataBytes)<dataSize || err != nil {
			// report error here
			break
		}

		// Create a decoder and receive a value.
		command := &Command{}
		payloadDecoderBuffer := bytes.NewBuffer(payloadBytes)
		payloadDecoder := gob.NewDecoder(payloadDecoderBuffer)

		err = payloadDecoder.Decode(command)
		if err != nil {
			log.Fatal("decode:", err)
		}
		command.Client = c

		fmt.Printf("Received %s on %s.%s %v\n", command.Action, command.Db, command.Table, command.Data)
		// deliver data to databases table
		databases[command.Db][command.Table].Chan<- command

		// push the line received above onto the channel
		//msgchan <- fmt.Sprintf("%s", line)
	}
}
