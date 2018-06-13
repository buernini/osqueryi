package osqueryi

import (
	"bytes"
	"encoding/binary"
	_ "fmt"
	"net"
	"os/exec"
)

// ++++++ Osquery server struct
type OsqServer struct {
	cmd     *bytes.Buffer
	History CmdHistoryInterface
}

// ++++++ Handler client request
func (q *OsqServer) handler(c *net.TCPConn) (res []byte, err error) {
	var run bool
	var n, length int
	var plen, buf []byte

	buf = make([]byte, 1024)
	q.cmd = bytes.NewBuffer([]byte{})
	for {
		if n, err = c.Read(buf); err != nil {
			return res, err
		}
		if n == 0 {
			continue
		}
		if run, err = q.writeCmd(buf); err != nil {
			return res, err
		}
		if run {
			if res, err = q.runCmd(); err != nil {
				return res, err
			} else {
				length = len(res)
				if length <= 4 {
					length = 0
				}
				plen = make([]byte, 8)
				binary.BigEndian.PutUint64(plen, uint64(length))
				c.Write(plen)
				if length > 0 {
					c.Write(res)
				}
				run = false
			}
			q.cmd.Reset()
		}
	}
	return res, nil
}

// ++++++ filter and write command to buffer
func (q *OsqServer) writeCmd(buf []byte) (run bool, err error) {
	for _, b := range buf {
		run = (b == 59)
		if b == 0 || run {
			break
		}
		if b == 10 {
			continue
		}
		if err = q.cmd.WriteByte(b); err != nil {
			return run, err
		}
	}
	return run, nil
}

// ++++++ Run command
func (q *OsqServer) runCmd() (b []byte, err error) {
	cmd := exec.Command("osqueryi", "--json", q.cmd.String())
	if b, err = cmd.Output(); err != nil {
		return bytes.TrimSpace(b), err
	}
	return bytes.TrimSpace(b), nil
}
