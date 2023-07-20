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

// The OpenAPI generated client functions generate a response with Body of []byte.
// This function formats the body into the required type.
func FormatResponseToType[T any](body []byte, newType *T) error {
	err := json.Unmarshal(body, &newType)
	if err != nil {
		return err
	}

	return nil
}
