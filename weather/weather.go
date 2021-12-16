package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type WeatherData struct {
	Description string `json:"description"`
}

type TempData struct {
	Temp float32 `json:"temp"`
}

type OpenWeatherApi struct {
	Location string        `json:"name"`
	Weather  []WeatherData `json:"weather"`
	Temp     TempData      `json:"main"`
	Code     string           `json:"cod"`
}

type OpenWeatherApiError struct {
	Message string `json:"message"`
	Code    string    `json:"cod"`
}

// Fetches weather from OpenWeather Api 
func FetchWeather(zipcode string) (OpenWeatherApi, error) {
	apiKey := os.Getenv("OPEN_WEATHER_API")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?zip=%v,us&units=imperial&appid=%v", zipcode, apiKey)

	weatherClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "weather-go")

	res, getErr := weatherClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		weather := OpenWeatherApi{}
		jsonErr := json.Unmarshal(body, &weather)

		if jsonErr != nil {
			fmt.Println(err)
		}

		fmt.Println(weather)
		return weather, nil

	} else {
		weatherError := OpenWeatherApiError{}
		jsonErr := json.Unmarshal(body, &weatherError)

		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		return OpenWeatherApi{}, errors.New(weatherError.Message)
	}

}