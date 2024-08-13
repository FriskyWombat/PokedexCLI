package pokeapi

import (
	"net/http"
	"github.com/FriskyWombat/pokedex/internal/jsonutil"
	//"github.com/FriskyWombat/pokedex/internal/pokecache"
)

const baseUrl = "https://pokeapi.co/api/v2/"

func getBaseUrl() string {
	return baseUrl
}

func GetFirstLocationUrl() string {
	return getBaseUrl() + "location?offset=0&limit=20"
}

func FetchData(url string, data interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	err = jsonutil.ParseResponse(res, &data)
	if err != nil {
		return err
	}
	return nil
}