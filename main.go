package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)


func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Weather Go: get your weather in JSON. Written in go")
	fmt.Println("Home Endpoint hit")
}

type WeatherData struct {
	Description string `json:"description"`
}

type TempData struct {
	Temp float32 `json:"temp"`
}

type OpenWeatherApi struct {
	Location string `json:"name"`
	Weather []WeatherData `json:"weather"`
	Temp TempData `json:"main"`
	Code int `json:"cod"`
}

type OpenWeatherApiError struct {
	Message string `json:"message"`
	Code int `json:"cod"`
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zip := vars["zip"]
	apiKey := os.Getenv("OPEN_WEATHER_API")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?zip=%v,us&units=imperial&appid=%v", zip, apiKey)

	weatherClient := http.Client {
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
			log.Fatal(jsonErr)
		}

		fmt.Println(weather)
		json.NewEncoder(w).Encode(weather)

	} else {
		weatherError := OpenWeatherApiError{}
		jsonErr := json.Unmarshal(body, &weatherError)

		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(weatherError)
	}


}

func handleReqs() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", home)
	myRouter.HandleFunc("/weather/{zip}", getWeather)
	log.Fatal(http.ListenAndServe(":5555", myRouter))
}



func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	handleReqs()
}