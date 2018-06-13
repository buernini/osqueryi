package osqueryi

import (
	"fmt"
	"net"
)

// +++++ Formater interface
type FormaterInterface interface {
	draw([]byte) []byte
}

// ++++++ Write command history
type CmdHistoryInterface interface {
	Write([]byte) (n int, err error)
	Read(p []byte) (n int, err error)
}

// ++++++ Run Server ...
func RunServer(host string, port int) error {
	var err error
	var network string = `tcp`
	var conn *net.TCPConn
	var tcpAddr *net.TCPAddr
	var listener *net.TCPListener
	var addr = fmt.Sprintf("%s:%d", host, port)

	if tcpAddr, err = net.ResolveTCPAddr(network, addr); err != nil {
		return err
	}
	if listener, err = net.ListenTCP(network, tcpAddr); err != nil {
		return err
	}

	for {
		if conn, err = listener.AcceptTCP(); err != nil {
			panic(err)
		}
		go func(c *net.TCPConn) {
			s := OsqServer{}
			s.handler(c)
		}(conn)
	}
	return nil
}

// ++++++ Run client...
func RunClient(host string, port int) error {
	client := &OsqClient{
		Host:     host,
		Port:     port,
		Formater: TableFormater{},
	}
	if err := client.run(); err != nil {
		return err
	}
	return nil
}
