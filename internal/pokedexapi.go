package internal

import (
	"encoding/json"
	"fmt"
	"github.com/mtslzr/pokeapi-go/structs"
	"io"
	"net/http"
)

func GetLocations(url string) ([]structs.Result, string, string, error) {
	var curl string
	var locations structs.Resource
	if url == "" {
		curl = "https://pokeapi.co/api/v2/location-area"
	} else {
		curl = url
	}
	fmt.Printf("curl: %s\n", curl)
	resp, err := http.Get(curl)
	if err != nil {
		return nil, "", "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", "", err
	}

	fmt.Printf("Response body: %s\n", body)
	var previousUrl string

	err = json.Unmarshal(body, &locations)
	if err != nil {
		return nil, "", "", err
	}
	if locations.Previous != nil {
		previousUrl = locations.Previous.(string)
	}

	return locations.Results, locations.Next, previousUrl, nil
}
