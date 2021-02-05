package handlers

import (
	"github.com/drbear95/gonotter-server/auth"
	"github.com/drbear95/gonotter-server/model"
	"github.com/graphql-go/graphql"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetNoteQueryResolver(p graphql.ResolveParams) (interface{}, error) {
	var result []model.Note
	var err error

	userId, err := auth.GetUserId(p.Context)

	if err != nil {
		return nil, err
	}

	search, searchOk := p.Args["search"].(string)
	take, takeOk := p.Args["take"].(int)
	skip, skipOk := p.Args["skip"].(int)

	opt := options.FindOptions{}

	query := append([]bson.M{}, bson.M{"author_id": bson.M{operator.Eq: userId}})

	if searchOk {
		query = append(query, bson.M{"$text": bson.M{"$search": search}})
	}

	if takeOk {
		take64 := int64(take)
		opt.Limit = &take64
	}

	if skipOk {
		skip64 := int64(skip)
		opt.Skip = &skip64
	}

	err = mgm.Coll(&model.Note{}).SimpleFind(&result, bson.M{"$and": query}, &opt)

	if err == nil && cap(result) > 0 {
		return result, nil
	}

	return nil, err
}

func GetCurrentUserQueryResolver(p graphql.ResolveParams) (interface{}, error) {
	var result []model.User
	var err error

	userId, err := auth.GetUserId(p.Context)

	if err != nil {
		return nil, err
	}

	err = mgm.Coll(&model.User{}).SimpleFind(&result, bson.M{"_id": bson.M{operator.Eq: userId}})

	if err == nil && cap(result) > 0 {
		return result, nil
	}

	return nil, err
}
