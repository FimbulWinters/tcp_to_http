package headers

import (
	"bytes"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	n = 0

	index := bytes.Index(data, []byte("\r\n"))
	if index == -1 {
		return n, false, nil //assuming that we havent been given enough data yet
	}
	if index == 0 {
		return n, true, nil
	}

}
