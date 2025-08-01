package geo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type GeoData struct {
	City string `json:"city"`
}

type CityPopulation struct {
	Error bool `json:"error"`
}

var ErrorNoCity = errors.New("NOCITY")
var ErrorNot200 = errors.New("NOT200")

func GetMyLocation(city string) (*GeoData, error) {
	if city != "" {
		isCity := CheckCity(city)
		if !isCity {
			return nil, ErrorNoCity
		}
		return &GeoData{
			City: city,
		}, nil
	}
	resp, err := http.Get("https://freegeoip.app/json/")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, ErrorNot200
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var geo GeoData
	json.Unmarshal(body, &geo)
	return &geo, nil
}

func CheckCity(city string) bool {
	postBody, _ := json.Marshal(map[string]string{
		"city": city,
	})
	resp, err := http.Post("https://countriesnow.space/api/v0.1/countries/population/cities/", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return false
	}
	defer resp.Body.Close() // дефер нужен, чтобы не произошла утечка памяти
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	var population CityPopulation
	json.Unmarshal(body, &population)
	return !population.Error
}
