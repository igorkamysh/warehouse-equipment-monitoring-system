package utils

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

func ParseRequestData[T any](reqBody io.Reader, p *T) error {
	dataBytes, err := io.ReadAll(reqBody)
	if err != nil {
		return errors.Wrap(err, "ParseRequestData: read body")
	}

	err = json.Unmarshal(dataBytes, p)
	if err != nil {
		return errors.Wrap(err, "ParseRequestData: unmarshal")
	}

	return nil
}
