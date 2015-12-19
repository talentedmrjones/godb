package engine

import (
	//"bufio"
	"fmt"
	//"io"
	"log"
	"net"
	"os"
)

/*
NewServer starts listening and initializes channels and kicks off various go routines
*/
func NewServer(databases Databases) {
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

		// run this in a goroutine so it doesnt block the for loop
		go handleConnection(conn, databases)
	} // end for
}

/*
handleConnection initializes new Connections
*/
func handleConnection(c net.Conn, databases Databases) {

	// when handleConnection finishes, execute c.Close to close the network connection
	defer c.Close()

	// initialize an instance of the Connection struct
	connection := NewConnection(c)

	// when the handleConnection finishes, execute this closure
	defer func() {
		// log to server console
		log.Printf("Connection from %s closed.\n", c.RemoteAddr())
		// push the connection into the rmchan for removal
		//rmchan <- connection
	}()

	// run this in a separate goroutine so as not to block connection.WriteLinesFrom
	go connection.send()
	connection.receive(databases)
}
