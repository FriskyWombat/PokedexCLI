package pokeapi

import (
	"net/http"
	"time"
	"fmt"
	"io"
	"github.com/FriskyWombat/pokedex/internal/jsonutil"
	"github.com/FriskyWombat/pokedex/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	cache pokecache.Cache
}
func NewClient() Client {
	return Client {
		httpClient: http.Client{
			Timeout: time.Second * 15,
		},
		cache: pokecache.NewCache(time.Second * 30),
	}
}

const baseUrl = "https://pokeapi.co/api/v2/"

func getBaseUrl() string {
	return baseUrl
}

func GetFirstLocationUrl() string {
	return getBaseUrl() + "location-area/?offset=0&limit=20"
}

func (c *Client) FetchData(url string, data interface{}) error {
	val, ok := c.cache.Get(url)
	if ok {
		err := jsonutil.ParseResponse(val, data)
		if err != nil {
			return err 
		}
		//fmt.Println("Cache hit on url:", url)
		return nil
	}
	res, err := c.FetchDataHTTP(url)
	if err != nil {
		return err
	}
	err = jsonutil.ParseResponse(res, data)
	if err != nil {
		return err 
	}
	//fmt.Println("Adding url to cache:", url)
	c.cache.Add(url, res)
	return nil
}

func (c *Client) FetchDataHTTP(url string) ([]byte, error) {
	res, err := c.httpClient.Get(url)
	defer res.Body.Close()
	if err != nil {
		return []byte{}, err
	}	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode > 299 {
		return []byte{}, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	return body, nil
}

