package main

import (
	"flag"
	"fmt"
	"osqueryi"
)

func main() {
	var host = flag.String(`h`, ``, `listen local ip address.`)
	var port = flag.Int(`p`, 0, `listen local port.`)
	flag.Parse()
	fmt.Printf("listen: %s:%d\n", *host, *port)
	if err := osqueryi.RunServer(*host, *port); err != nil {
		panic(err)
	}
}
