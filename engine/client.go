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
	socket		net.Conn 		// holds the network socket
	replies 	chan *Reply	// receives replies
}

// Receive continuously looks for data from the socket and relays that to a table's command channel
func (c *Client) Receive(databases map[string]map[string]*Table) {
	// create a buffered reader
	buf := bufio.NewReader(c.socket)

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
		command.client = c

		fmt.Printf("Received %s on %s.%s %v\n", command.Action, command.Db, command.Table, command.Data)
		// deliver data to databases table
		databases[command.Db][command.Table].commands<- command

		// push the line received above onto the channel
		//msgchan <- fmt.Sprintf("%s", line)
	}
}


// ReadLinesInto continuously looks for data from the connection and relays that to the message channel
func (c *Client) Send() {
	for reply := range c.replies {
		fmt.Printf("reply: %v\n", reply)
		// gob encode reply into payload
		var payloadEncodingBuffer bytes.Buffer
		payloadEncoder := gob.NewEncoder(&payloadEncodingBuffer)
		payloadEncodingErr := payloadEncoder.Encode(reply)
		if payloadEncodingErr != nil {
			// TODO handle error
			fmt.Printf("payloadEncodingErr %v\n", payloadEncodingErr)
		}
		payloadBytes := payloadEncodingBuffer.Bytes()

		//fmt.Printf("%v %v", encodeErr, payloadBytes)
		dataSize := make([]byte,4)
		binary.BigEndian.PutUint32(dataSize, uint32(len(payloadBytes)))

		c.socket.Write(dataSize)
		c.socket.Write(payloadBytes)
	}
}
