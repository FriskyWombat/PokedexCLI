package pokeapi

import (
	"net/http"
	"github.com/FriskyWombat/pokedex/internal/jsonutil"
)

const baseUrl = "https://pokeapi.co/api/v2/"

func GetBaseUrl() string {
	return baseUrl
}

func GetFirstLocationUrl() string {
	return GetBaseUrl() + "location?offset=0&limit=20"
}

type LocationResp struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func FetchMapData(url string) (LocationResp, error) {
	res, err := http.Get(url)
	if err != nil {
		return LocationResp{}, err
	}
	data := LocationResp{}
	err = jsonutil.ParseResponse(res, &data)
	if err != nil {
		return LocationResp{}, err
	}
	return data, nil
}