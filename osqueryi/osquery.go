package osqueryi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	typeUnknow      uint16 = 0
	typeError       uint16 = 1
	typeResNormal   uint16 = 2
	typeReqQueryCmd uint16 = 3
	typeReqAuth     uint16 = 4
	typeReqHelpCmd  uint16 = 5
)

const (
	HeaderLength int = 10
)

// ++++++ Server protocol header
type ProtocolHeader struct {
	Type          uint16
	ContentLength uint64
}

// ++++++ Pack
func (self ProtocolHeader) pack(pType uint16, pLength uint64) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, pType)
	binary.Write(buf, binary.BigEndian, pLength)
	return buf.Bytes()
}

// +++++ Formater interface
type FormaterInterface interface {
	draw([]byte) []byte
}

// ++++++ Run Server ...
func RunServer(host string, port int, path string) error {
	var err error
	var network string = `tcp`
	var conn *net.TCPConn
	var tcpAddr *net.TCPAddr
	var listener *net.TCPListener
	var addr = fmt.Sprintf("%s:%d", host, port)
	session := new(OsqSession)
	if err = session.init(path); err != nil {
		return err
	}

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
		go func(c *net.TCPConn, session *OsqSession) {
			s := OsqServer{session: session}
			s.handler(c)
		}(conn, session)
	}
	return nil
}

// ++++++ Run client...
func RunClient(host string, port int, user, password string) error {
	client := &OsqClient{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Formater: TableFormater{},
	}
	if err := client.run(); err != nil {
		return err
	}
	return nil
}

type Record struct {
	c *net.TCPConn
}

func (self *Record) read() (reqType uint16, data []byte, err error) {
	var packHeader ProtocolHeader
	header := make([]byte, HeaderLength)
	if _, err = self.c.Read(header); err != nil {
		return typeUnknow, data, err
	}
	if err = binary.Read(bytes.NewReader(header), binary.BigEndian, &packHeader); err != nil {
		return typeUnknow, data, err
	}
	data = make([]byte, packHeader.ContentLength)
	if _, err = self.c.Read(data); err != nil {
		return typeUnknow, data, err
	}
	return packHeader.Type, data, nil
}

func (self *Record) write(resType uint16, data []byte) (err error) {
	var packHeader ProtocolHeader
	h := packHeader.pack(resType, uint64(len(data)))
	if _, err = self.c.Write(h); err != nil {
		return err
	}
	if _, err = self.c.Write(data); err != nil {
		return err
	}
	return nil

}
