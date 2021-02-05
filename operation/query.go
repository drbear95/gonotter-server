package operation

import (
	"github.com/drbear95/gonotter-server/ghql/concretes"
	"github.com/drbear95/gonotter-server/model"
	"github.com/graphql-go/graphql"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
)

var Query = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type:        concretes.UserType,
				Description: "Get users",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					var result []model.User
					var err error
					if ok {
						err = mgm.Coll(&model.User{}).SimpleFind(&result, bson.M{"id": bson.M{operator.Eq: id}})
						if err != nil && cap(result) > 0{
							return result, nil
						}
					}else{
						err = mgm.Coll(&model.User{}).SimpleFind(&result, bson.M{})
						if err != nil && cap(result) > 0 {
							return result, nil
						}
					}
					return nil, err
				},
			},
			//"note": &ghql.Field{
			//	Type:        model.NoteType,
			//	Description: "Get notes",
			//	Args: ghql.FieldConfigArgument{
			//		"id": &ghql.ArgumentConfig{
			//			Type: ghql.String,
			//		},
			//		"name": &ghql.ArgumentConfig{
			//			Type: ghql.String,
			//		},
			//		"authorId": &ghql.ArgumentConfig{
			//			Type: ghql.String,
			//		},
			//	},
			//	Resolve: func(p ghql.ResolveParams) (interface{}, error) {
			//		id, isIdOk := p.Args["id"].(string)
			//		id, isIdOk := p.Args["name"].(string)
			//		id, isIdOk := p.Args["authorId"].(string)
			//		var result []model.User
			//		var err error
			//		if ok {
			//			err = mgm.Coll(&model.User{}).SimpleFind(result, bson.M{"id": bson.M{operator.Eq: primitive.ObjectIDFromHex(id)}})
			//			if err != nil {
			//				return result, nil
			//			}
			//		}
			//		return nil, err
			//	},
			//},
		},
	})
