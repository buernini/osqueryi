package main

import (
	"flag"
	"fmt"
	"osqueryi"
)

func main() {

	/*
		auth := &osqueryi.Auth{}
		auth.Load(`/usr/local/golang/osqueryi/src/config.txt`)
		err := auth.Login(`wanglipeng`, `111111`)
		fmt.Println(err)
		return
	*/

	var host = flag.String(`h`, ``, `listen local ip address.`)
	var port = flag.Int(`P`, 0, `listen local port.`)
	var conf = flag.String(`c`, ``, `login account config file.`)
	flag.Parse()
	fmt.Printf("listen: %s:%d\n", *host, *port)
	if err := osqueryi.RunServer(*host, *port, *conf); err != nil {
		panic(err)
	}
}
