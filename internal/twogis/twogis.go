package twogis

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Meta struct {
	ApiVersion string `json:"api_version"`
	Code       int    `json:"code"`
	IssueDate  string `json:"issue_date"`
}

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type ResponseItems struct {
	AddressName string `json:"address_name"`
	FullName    string `json:"full_name"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Point       Point  `json:"point"`
	PurposeName string `json:"purpose_name"`
	Type        string `json:"type"`
}

type Result struct {
	Items []ResponseItems `json:"items"`
	Total int             `json:"total"`
}

type Response struct {
	Meta   Meta   `json:"meta"`
	Result Result `json:"result"`
}

func GetLatLon(address string) (float64, float64, error) {
	if address == "" {
		return 0, 0, errors.New("empty address")
	}
	twoGisApiKey := os.Getenv("API_KEY_2GIS")
	queryUrl := "https://catalog.api.2gis.com/3.0/items/geocode"

	client := &http.Client{}
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return 0, 0, err
	}
	if req != nil {
		q := req.URL.Query()
		q.Add("q", address)
		q.Add("fields", "items.point")
		q.Add("key", twoGisApiKey)
		req.URL.RawQuery = q.Encode()
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer res.Body.Close()

	if res.StatusCode > 200 {
		return 0, 0, errors.New("status code: " + strconv.Itoa(res.StatusCode))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, 0, err
	}
	var response Response
	err = json.Unmarshal(body, &response)

	if response.Meta.Code > 299 {
		return 0, 0, errors.New("2gis non 2xx response code")
	}

	if response.Result.Total == 0 {
		return 0, 0, errors.New("place was not found")
	}
	// What should I do if total count greater than 1
	point := response.Result.Items[0].Point
	if point.Lat == 0 || point.Lon == 0 {
		return 0, 0, errors.New("point not found")
	}
	return point.Lat, point.Lon, nil
}
