package handlers

import (
	"errors"
	"github.com/drbear95/gonotter-server/auth"
	"github.com/drbear95/gonotter-server/model"
	"github.com/graphql-go/graphql"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetNoteMutationResolver(p graphql.ResolveParams) (interface{}, error) {
	id, idOk := p.Args["id"].(string)
	title, titleOk := p.Args["title"].(string)
	content, contentOk := p.Args["content"].(string)

	var err error

	userId, err := auth.GetUserId(p.Context)

	if err != nil {
		return nil, err
	}

	if !titleOk {
		return nil, errors.New("title is not ok")
	}

	if !contentOk {
		return nil, errors.New("content is not ok")
	}

	note := model.NewNote(title, content, *userId)

	if idOk {
		err = updateNote(id, note)
	} else {
		err = createNote(note)
	}

	err = mgm.Coll(note).Create(note)

	if err != nil {
		return note, nil
	} else {
		return nil, err
	}
}

func createNote(note *model.Note) error {
	return mgm.Coll(note).Create(note)
}

func updateNote(id string, note *model.Note) error {
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	note.ID = objId

	return mgm.Coll(note).Update(note)
}

