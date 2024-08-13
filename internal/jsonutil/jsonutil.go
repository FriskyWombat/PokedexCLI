package jsonutil

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
)


func ParseResponse(res *http.Response, v interface{}) error {
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &v)
	if err != nil {
		return err
	}
	return nil
}