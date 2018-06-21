package osqueryi

import (
	"fmt"
	"net"
	"strings"
)

// ++++++ Osquery server struct
type OsqServer struct {
	session *OsqSession
}

// ++++++ Handler client request
func (self *OsqServer) handler(c *net.TCPConn) (err error) {
	var record Record
	var reqType uint16
	var clientId net.Addr
	var res []byte
	var cmd []byte

	cmder := new(OsqCommand)
	cmder.init()
	clientId = c.RemoteAddr()
	record = Record{c: c}

	for {
		if reqType, cmd, err = record.read(); err != nil {
			break
		}
		cmder.cmd.Write(cmd)

		// process by header type
		switch reqType {
		case typeReqAuth:
			auth := strings.Split(cmder.String(), ` `)
			if len(auth) != 2 {
				record.write(typeError, []byte(`user or password miss`))
				continue
			}
			if err = self.session.login(clientId.String(), auth[0], auth[1]); err != nil {
				record.write(typeError, []byte(`login failed`))
				continue
			}
		default:
			if !self.session.valid(clientId.String()) {
				record.write(typeError, []byte(`session invald`))
				continue
			}
			res, err = cmder.query()
			if err != nil {
				record.write(typeError, []byte(fmt.Sprintf("%s", err)))
				cmder.reset()
				continue
			}
			record.write(typeResNormal, res)
		}
		cmder.reset()
	}
	return err
}
