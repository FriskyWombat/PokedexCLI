package jsonutil

import (
	"encoding/json"
)


func ParseResponse(body []byte, v interface{}) error {
	err := json.Unmarshal(body, &v)
	if err != nil {
		return err
	}
	return nil
}