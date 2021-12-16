package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/opti21/weather-go/graph"
	"github.com/opti21/weather-go/graph/generated"
	"github.com/opti21/weather-go/weather"
)

const defaultPort = "8080"

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Weather Go: get your weather in JSON and Grahql. Written in go")
	fmt.Println("Home Endpoint hit")
}

// Parses zipcode from url request to send to FetchWeather function and returns JSON
func getWeather(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zip := vars["zip"]

	fetchedWeather, weatherErr := weather.FetchWeather(zip)

	if  weatherErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(weatherErr)
	} else {
		fmt.Println("Got weather")
		fmt.Println(fetchedWeather)
		json.NewEncoder(w).Encode(fetchedWeather)

	}


}

func handleReqs() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", home)
	myRouter.HandleFunc("/rest/{zip}", getWeather)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	// myRouter.Handle("/", playground.Handler("GraphQL playground", "/query"))
	myRouter.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, myRouter))
}

func main() {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	handleReqs()
}