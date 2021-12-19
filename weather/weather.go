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
	Temp float64 `json:"temp"`
}

type OpenWeatherApiResponse struct {
	Location string        `json:"name"`
	Weather  []WeatherData `json:"weather"`
	Temp     TempData      `json:"main"`
	Message  string        `json:"message,omitempty"`
	Code     float64       `json:"cod,string"`
}

type Weather struct {
	Location         string
	CurrentCondition string
	Temp             float64
	Code             float64
	Message          string
}

type RequestError struct {
	StatusCode float64
	Err        error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("OpenWeather Error: status %v: err %q", r.StatusCode, r.Err)
}

// Fetches weather from OpenWeather Api
func FetchWeather(zipcode string) (Weather, error) {
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

	apiResponse := OpenWeatherApiResponse{}
	jsonErr := json.Unmarshal(body, &apiResponse)

	if jsonErr != nil {
		fmt.Println(err)
	}

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		weather := Weather{
			Location:         apiResponse.Location,
			CurrentCondition: apiResponse.Weather[0].Description,
			Temp:             apiResponse.Temp.Temp,
			Code:             apiResponse.Code,
		}

		fmt.Println(weather)
		return weather, nil

	} else {
		return Weather{}, &RequestError{
			StatusCode: apiResponse.Code,
			Err:        errors.New(apiResponse.Message),
		}
	}

}
