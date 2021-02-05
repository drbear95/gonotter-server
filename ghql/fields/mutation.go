package fields

import (
	"github.com/drbear95/gonotter-server/ghql/concretes"
	"github.com/drbear95/gonotter-server/ghql/handlers"
	"github.com/graphql-go/graphql"
)

func GetNoteMutationField() *graphql.Field {
	return &graphql.Field{
		Type:        concretes.NoteType,
		Description: "Notes mutation",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"title": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"content": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: handlers.GetNoteMutationResolver,
	}
}