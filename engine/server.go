package engine

import (
	//"bufio"
	"fmt"
	//"io"
	"log"
	"net"
	"os"

)

// RunServer starts listening and initializes channels and kicks off various go routines
func RunServer(databases map[string]map[string]*Table) {
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// loop forever waiting on new connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// run this in a thread so it doesnt block the for loop
		go handleConnection(conn, databases)
	} // end for
}


// handleConnection initializes new Clients
func handleConnection(c net.Conn, databases map[string]map[string]*Table) {

	// when handleConnection finishes, execute c.Close thereby closing the network connection
	defer c.Close()

	// initialize an instance of the Client struct
	client := Client{
		conn:     c, // store the network connection
		//ch:       make(chan string), // store a channel for strings
	}

	// when the handleConnection finishes, execute this closure
	defer func() {
		// log to server console
		log.Printf("Connection from %s closed.\n", c.RemoteAddr())
		// push the client into the rmchan for removal
		//rmchan <- client
	}()

	// run this in a separate thread so as not to block client.WriteLinesFrom
	//go client.Send()
	client.Receive(databases)
}
