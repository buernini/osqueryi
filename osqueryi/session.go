package osqueryi

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"sync"
	"time"
)

// ++++++ Session
type OsqSession struct {
	mux     sync.Mutex
	clients map[string]time.Time
	account map[string]string
}

// ++++++ Load accounts
func (self *OsqSession) init(path string) (err error) {
	self.mux.Lock()
	defer self.mux.Unlock()
	self.clients = make(map[string]time.Time)
	self.account = make(map[string]string)

	var fd *os.File
	if fd, err = os.Open(path); err != nil {
		return err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		items := bytes.Split([]byte(scanner.Text()), []byte(`=`))
		if len(items) != 2 {
			break
		}
		self.account[string(items[0])] = string(items[1])
	}
	return nil
}

// ++++++ Login
func (self *OsqSession) login(uid, user, password string) (err error) {
	if pwd, ok := self.account[user]; !ok || pwd != password {
		return errors.New(`login failed`)
	}
	self.clients[uid] = time.Now()
	return nil
}

// ++++++ Valid session status
func (self *OsqSession) valid(uid string) bool {
	var ok bool
	var loginTime time.Time
	loginTime, ok = self.clients[uid]
	return ok && loginTime.Unix() > 0
}
