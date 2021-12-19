package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/opti21/weather-go/graph"
	"github.com/opti21/weather-go/graph/generated"
	"github.com/opti21/weather-go/weather"
)

const defaultPort = "8080"

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<h1>Weather Go</h1>
	<p>Get your weather in JSON <b>(/rest/{zip})</b> and Grahql <b>(/query)</b>.</p>
	<p>Data provided by <a href="https://openweathermap.org/" target="_blank">OpenWeather</a>
	<p>Written in go</p>
	<h2>GraphQL example</h2>
	<pre>
		<code>
		query {
		  currentWeather(zipcode: "77001") {
		    condition
		    location
		    zipcode
		    temp
		  }
		}
		</code>
	</pre>
	`)
	fmt.Println("Home Endpoint hit")
}

// Parses zipcode from url request to send to FetchWeather function and returns JSON
func getWeather(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zip := vars["zip"]

	fetchedWeather, weatherErr := weather.FetchWeather(zip)

	if  weatherErr != nil {
		wthErr := struct{
			Error string
			Code  float64
		}{
			Error: weatherErr.Error(),
			Code: 500,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(wthErr)
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

	myRouter.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	myRouter.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, myRouter))
}

func main() {
	handleReqs()
}