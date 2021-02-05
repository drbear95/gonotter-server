package ghql

import (
	"github.com/drbear95/gonotter-server/ghql/fields"
	"github.com/graphql-go/graphql"
	"log"
)

func GetSchema() *graphql.Schema{
	graphqlSchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"note": fields.GetNoteQueryField(),
				"currentUser": fields.GetCurrentUserQueryField(),
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"newNote": fields.GetNoteMutationField(),
			},
		}),
	})

	if err != nil {
		log.Fatal(err)
	}

	return &graphqlSchema
}