package main

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
)

type FetchedData struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func parseResponse(res *http.Response, v interface{}) error {
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