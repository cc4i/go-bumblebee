package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"gql/airclt"
	"gql/graph/generated"
	"gql/graph/model"

	log "github.com/sirupsen/logrus"
)

func (r *mutationResolver) Save(ctx context.Context, input model.NewAirQuality) (*model.AirQuality, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AirQuality(ctx context.Context, city string) (*model.AirQuality, error) {
	var air model.AirQuality
	buf, err := airclt.AirOfCity(ctx, city)
	if err != nil {
		log.Error(err.Error())
		return &air, err
	}
	if err = json.Unmarshal(buf, &air); err != nil {
		log.WithField("from", "air-service").Errorf("%s", buf)
		log.Error(err.Error())
		return &air, err
	}
	return &air, nil
}

func (r *queryResolver) AirQualityHistory(ctx context.Context, city string) ([]*model.AirQuality, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
