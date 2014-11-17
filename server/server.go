package server

import (
	"fmt"
	"bufio"
	"log"
	"net"
)

func Run() {

	// Listen on TCP port 2000 on all interfaces.
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			connReader := bufio.NewReader(conn)
	    for {
	        line, err := connReader.ReadBytes(byte("\u2404"))
					if err!=nil {
						panic(err)
					}
					fmt.Printf("%v", line)
	    }
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}
