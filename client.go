package main

import (
	"flag"
	"osqueryi"
)

func main() {
	var host = flag.String(`h`, ``, `remote server ip address.`)
	var port = flag.Int(`p`, 0, `remote server port.`)
	flag.Parse()
	if err := osqueryi.RunClient(*host, *port); err != nil {
		panic(err)
	}
}
