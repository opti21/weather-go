package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/opti21/weather-go/graph/generated"
	"github.com/opti21/weather-go/graph/model"
	"github.com/opti21/weather-go/weather"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *queryResolver) CurrentWeather(ctx context.Context, zipcode string) (*model.CurrentWeather, error) {
	fetchedWeather, err := weather.FetchWeather(zipcode)

	if err != nil {
		graphql.AddError(ctx, err)
		return nil, gqlerror.Errorf("API Error")
	}


	current := model.CurrentWeather{
		Location:  fetchedWeather.Location,
		Condition: fetchedWeather.Weather[0].Description,
		Zipcode:   zipcode,
		Temp:      float64(fetchedWeather.Temp.Temp),
	}
	return &current, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
