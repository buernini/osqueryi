package main

import (
	"flag"
	"osqueryi"
)

func main() {
	var host = flag.String(`h`, ``, `remote server ip address.`)
	var port = flag.Int(`P`, 0, `remote server port.`)
	var user = flag.String(`u`, ``, `username.`)
	var password = flag.String(`p`, ``, `password.`)
	flag.Parse()
	if err := osqueryi.RunClient(*host, *port, *user, *password); err != nil {
		//panic(err)
	}
}
