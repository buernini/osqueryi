package osqueryi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ++++++ Fields sort
type Fields []string

func (f Fields) Len() int {
	return len(f)
}

func (f Fields) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f Fields) Less(i, j int) bool {
	return (strings.Compare(f[i], f[j]) == -1)
}

// +++++ Default
type DefaultFormater struct {
}

// ++++++ draw
func (f DefaultFormater) draw(data []byte) []byte {
	var dataItems []map[string]string
	json.Unmarshal(data, &dataItems)

	var buf *bytes.Buffer
	buf = bytes.NewBuffer([]byte{})
	for _, values := range dataItems {
		for k, v := range values {
			buf.WriteString(fmt.Sprintf("%s:\t%s\n", k, v))
		}
		buf.WriteString("\n")
	}
	buf.WriteString(fmt.Sprintf("total:%d\n", len(dataItems)))
	return buf.Bytes()
}

// ++++++ Table
type TableFormater struct {
	fieldSpace map[string]int
}

func (f TableFormater) draw(data []byte) []byte {
	var dataItems []map[string]string
	if err := json.Unmarshal(data, &dataItems); err != nil {
		return []byte{}
	}

	if f.fieldSpace == nil {
		f.fieldSpace = make(map[string]int)
	}

	var fields Fields
	fields = make([]string, len(dataItems[0]))
	var i int
	for k, _ := range dataItems[0] {
		if k == `` {
			continue
		}
		if len(k) > f.fieldSpace[k] {
			f.fieldSpace[k] = len(k)
		}
		fields[i] = k
		i++
	}
	sort.Sort(fields)

	var buf *bytes.Buffer
	buf = bytes.NewBuffer([]byte{})
	for _, values := range dataItems {
		for _, field := range fields {
			value := values[field]
			if len(value) > f.fieldSpace[field] {
				f.fieldSpace[field] = len(value)
			}
		}
	}

	var width int
	for _, v := range fields {
		space := f.fieldSpace[v]
		width += space
		sep := f.separatorBySpace(space-len(v), ' ')
		buf.WriteString(fmt.Sprintf(" %s%s |", v, string(sep)))
	}
	width += (len(fields) * 3)
	sep := f.separatorBySpace(width, '_')
	buf.WriteString(fmt.Sprintf("\n%s\n\n", string(sep)))

	for _, values := range dataItems {
		for _, field := range fields {
			value := values[field]
			space := f.fieldSpace[field]
			sep := f.separatorBySpace(space-len(value), ' ')
			buf.WriteString(fmt.Sprintf(" %s%s |", value, string(sep)))
		}
		buf.WriteString("\n")
	}
	sep = f.separatorBySpace(width, '_')
	buf.WriteString(fmt.Sprintf("%s\n\n", string(sep)))

	return buf.Bytes()
}

func (f TableFormater) separatorBySpace(space int, b byte) (sep []byte) {
	sep = make([]byte, space)
	for i := 0; i < space; i++ {
		sep[i] = b
	}
	return sep
}

// ++++++ Other
type OtherFormater struct{}

func (f OtherFormater) draw(data []byte) []byte {
	// to do
	return data
}
