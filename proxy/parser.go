package proxy

import (
	"encoding/json"
)

func Parse(data []byte) (*ConnectionMeta, error) {
	message := &ConnectionMeta{}
	err := json.Unmarshal(data, &message)

	if err == nil {
		return message, err
	}
	return nil, err
}
