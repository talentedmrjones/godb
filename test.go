package main

import (
	"fmt"
	"time"
)

func main () {
	now := time.Now()
	fmt.Printf("%v\n", time.Since(now))
	fmt.Printf("%v\n", time.Since(now))
	fmt.Printf("%v\n", time.Since(now))
	fmt.Printf("%v\n", time.Since(now))
	fmt.Printf("%v\n", time.Since(now))
}
