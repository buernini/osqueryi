package osqueryi

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"
)

// ++++++ Osquery client struct
type OsqClient struct {
	conn           net.Conn
	Host           string
	Port           int
	User, Password string
	Formater       FormaterInterface
}

// ++++++ Build connect
func (self *OsqClient) connect() (err error) {
	var addr string = fmt.Sprintf("%s:%d", self.Host, self.Port)
	var network string = `tcp`
	var dialer *net.Dialer

	dialer = &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 3 * time.Second,
	}

	if self.conn, err = dialer.Dial(network, addr); err != nil {
		return err
	}
	return nil
}

// ++++++ Client run
func (self *OsqClient) run() (err error) {
	var c byte
	var response []byte
	var resType uint16

	if err = self.connect(); err != nil {
		return err
	}
	defer self.conn.Close()

	cmder := new(OsqCommand)
	cmder.init()

	// send auth data
	authData := []byte(fmt.Sprintf("%s %s", self.User, self.Password))
	if err = self.execCmd(typeReqAuth, authData); err != nil {
		return err
	}
	self.toStdout([]byte{})

	if self.Formater == nil {
		self.Formater = DefaultFormater{}
	}

	for {
		fmt.Scanf("%c", &c)
		if done, _ := self.cmdAnalyze(cmder, c); !done {
			continue
		}

		/*
			if quit := self.isQuit(); quit {
				break
			}
		*/
		if err = self.execCmd(cmder.Type, cmder.cmd.Bytes()); err != nil {
			self.toStdout([]byte(fmt.Sprintf("%s", err)))
			continue
		}
		cmder.reset()

		if resType, response, err = self.result(); err != nil {
			self.toStdout([]byte(fmt.Sprintf("%s", err)))
			continue
		}

		if resType == typeError {
			self.toStdout(response)
			continue
		}
		if resType != typeResNormal {
			self.toStdout([]byte(`unknow resType`))
			continue
		}

		switch cmder.Type {
		case typeReqHelpCmd:
			self.toStdout(response)
		default:
			self.toStdout(self.Formater.draw(response))
		}
	}
	return err
}

func (self *OsqClient) cmdAnalyze(cmder *OsqCommand, c byte) (bool, error) {
	if len(cmder.cmd.Bytes()) == 0 {
		if c == 46 {
			cmder.Type = typeReqHelpCmd
		} else {
			cmder.Type = typeReqQueryCmd
		}
	}
	if cmder.Type == typeReqHelpCmd && c == 10 {
		return true, nil
	}
	if cmder.Type == typeReqQueryCmd && c == 59 {
		return true, nil
	}

	if c == 10 {
		return false, nil
	}

	return false, cmder.cmd.WriteByte(c)
}

// ++++++ Check client quit command
func (self *OsqClient) isQuit() (quit bool) {
	buf := []byte{}
	if len(buf) != 5 {
		return false
	}
	if quit = (bytes.Compare(buf[0:5], []byte(`quit;`)) == 0); quit {
	}
	return quit
}

func (self *OsqClient) execCmd(reqType uint16, cmd []byte) (err error) {
	var record Record
	record = Record{c: self.conn.(*net.TCPConn)}
	return record.write(reqType, cmd)
}

func (self *OsqClient) result() (resType uint16, data []byte, err error) {
	var record Record
	record = Record{c: self.conn.(*net.TCPConn)}
	return record.read()
}

func (t *OsqClient) toStdout(data []byte) {
	fmt.Fprintf(os.Stdout, "%s\n", data)
	fmt.Fprintf(os.Stdout, "\033[31m%s:%d> \033[0m", t.Host, t.Port)
}
