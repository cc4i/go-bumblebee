package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"gql/airclt"
	"gql/graph/generated"
	"gql/graph/model"

	log "github.com/sirupsen/logrus"
)

func (r *mutationResolver) SaveQueryHistory(ctx context.Context, input model.NewQuery) (*model.QueryHistory, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AirQuality(ctx context.Context, city string) (*model.AirQuality, error) {
	air, err := airclt.AirOfCity(ctx, city)
	if err != nil {
		log.WithField("from", "air-service").Error(err)
		return air, err
	}

	return air, nil
}

func (r *queryResolver) QueryHistory(ctx context.Context, city string) ([]*model.QueryHistory, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
