package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// WeatherInfo contains the weather information
type WeatherInfo struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		Pressure  float64 `json:"pressure"`
		Humidity  int     `json:"humidity"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		SeaLevel  float64 `json:"sea_level"`
		GrndLevel float64 `json:"grnd_level"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}


func getWeatherData(token string, lat string, lon string) (*WeatherInfo, error) {
	var weatherinfo WeatherInfo

	baseURL, err := url.Parse("https://api.openweathermap.org/data/2.5/weather?")
	if err != nil {
		fmt.Println("Malformed URL: ", err.Error())
	}

	params := url.Values{}
	params.Add("lat", lat)
	params.Add("lon", lon)
	params.Add("appid", token)


	baseURL.RawQuery = params.Encode()

	fmt.Printf("Encoded URL is %q\n", baseURL.String())

	response, err := http.Get(baseURL.String())
	if err != nil {
		return &weatherinfo, err
	}

	data, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal([]byte(data), &weatherinfo)

	return &weatherinfo, nil
}


func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("service accessed")

	lons, ok := r.URL.Query()["lon"]

	if !ok || len(lons[0]) < 1 {
		log.Println("Url Param 'lon' is missing")
		return
	}

	lats, ok := r.URL.Query()["lat"]

	if !ok || len(lats[0]) < 1 {
		log.Println("Url Param 'lat' is missing")
		return
	}

	appids, ok := r.URL.Query()["appid"]

	if !ok || len(appids[0]) < 1 {
		log.Println("Url Param 'appid' is missing")
		return
	}

	lat := lats[0]
	lon := lons[0]
	appid := appids[0]

	fmt.Println("Longitude is ", lon)
	fmt.Println("Latitude is ", lat)
	fmt.Println("Appid is ", appid)

	weatherinfo, err := getWeatherData(appid, lat, lon)
	if err != nil {
		fmt.Println("Error occurred ", err)
	}

	js, err := json.Marshal(weatherinfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Write(js)
}

func main() {

	http.HandleFunc("/climate", handler)
	http.ListenAndServe(":3000", nil)
}

