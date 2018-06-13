package osqueryi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

// ++++++ Osquery client struct
type OsqClient struct {
	conn     net.Conn
	Host     string
	Port     int
	Formater FormaterInterface
	History  CmdHistoryInterface
}

// ++++++ Build connect
func (t *OsqClient) connect() (err error) {
	var addr string = fmt.Sprintf("%s:%d", t.Host, t.Port)
	var network string = `tcp`
	var dialer *net.Dialer

	dialer = &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 3 * time.Second,
	}

	if t.conn, err = dialer.Dial(network, addr); err != nil {
		return err
	}
	return nil
}

// ++++++ Client run
func (t *OsqClient) run() (err error) {
	var c byte
	var buf *bytes.Buffer
	var response []byte

	buf = bytes.NewBuffer([]byte{})
	if err = t.connect(); err != nil {
		return err
	}
	t.guide()
	if t.Formater == nil {
		t.Formater = DefaultFormater{}
	}
	for {
		fmt.Scanf("%c", &c)
		if c == 10 {
			continue
		}
		buf.WriteByte(c)
		if c == 59 {
			if quit := t.isQuit(buf.Bytes()); quit {
				buf.Reset()
				t.conn.Close()
				break
			}
			if response, err = t.exec(buf.Bytes()); err != nil {
				buf.Reset()
				break
			}
			buf.Reset()
			fmt.Fprintf(os.Stdout, "%s", t.Formater.draw(response))
			t.guide()
		}
	}
	return err
}

// ++++++ Check client quit command
func (t *OsqClient) isQuit(cmd []byte) bool {
	return bytes.Compare(cmd[0:5], []byte(`quit;`)) == 0
}

// ++++++ Write client command to server
func (t *OsqClient) exec(cmd []byte) (response []byte, err error) {
	var n int
	var proLen []byte
	var resLen uint64

	if _, err = t.conn.Write(cmd); err != nil {
		return response, err
	}
	// todo : save cmd history to disk

	proLen = make([]byte, 8)
	if _, err = t.conn.Read(proLen); err != nil {
		return response, err
	}
	resLen = binary.BigEndian.Uint64(proLen)
	for resLen > 0 {
		buf := make([]byte, 1024)
		if n, err = t.conn.Read(buf); err != nil {
			return response, err
		}
		response = append(response, buf[0:n]...)
		resLen -= uint64(n)
	}

	return response, nil
}

func (t *OsqClient) guide() {
	fmt.Fprintf(os.Stdout, "\033[31m%s:%d> \033[0m", t.Host, t.Port)
}
