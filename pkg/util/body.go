package util

import (
	"bytes"
	"encoding/json"
	"io"
)

func MakeBodyReaderFromType[T any](body T) (io.Reader, error) {
	var bodyReader io.Reader

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader = bytes.NewReader(buf)
	return bodyReader, nil
}
