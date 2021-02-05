package fields

import (
	"github.com/drbear95/gonotter-server/ghql/concretes"
	"github.com/drbear95/gonotter-server/ghql/handlers"
	"github.com/graphql-go/graphql"
)

func GetNoteQueryField() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(concretes.NoteType),
		Description: "Notes query",
		Args: graphql.FieldConfigArgument{
			"search": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"take": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"skip": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: handlers.GetNoteQueryResolver,
	}
}

func GetCurrentUserQueryField() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(concretes.UserType),
		Description: "Current user query",
		Resolve: handlers.GetCurrentUserQueryResolver,
	}
}