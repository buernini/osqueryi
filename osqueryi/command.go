package osqueryi

import (
	"bytes"
	"errors"
	"os/exec"
)

type OsqCommand struct {
	Type uint16
	cmd  *bytes.Buffer
}

func (self *OsqCommand) init() {
	self.cmd = bytes.NewBuffer([]byte{})
}

func (self *OsqCommand) query() (b []byte, err error) {
	var stderr bytes.Buffer

	// trim end mark
	//self.cmd = bytes.TrimSuffix(self.cmd, []byte(`;`))

	cmd := exec.Command("osqueryi", "--json", self.String())
	cmd.Stderr = &stderr
	if b, err = cmd.Output(); err != nil {
		return bytes.TrimSpace(b), err
	}
	if len(stderr.Bytes()) > 0 {
		err = errors.New(string(stderr.Bytes()))
	}
	return bytes.TrimSpace(b), err
}

func (self *OsqCommand) write(c byte) error {
	return self.cmd.WriteByte(c)
}

func (self *OsqCommand) length() int {
	return len(self.cmd.Bytes())
}

func (self *OsqCommand) reset() {
	self.cmd.Reset()
}

func (self *OsqCommand) String() string {
	return string(self.cmd.Bytes())
}
