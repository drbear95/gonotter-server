package persistance

import (
	"github.com/drbear95/gonotter-server/model"
	"github.com/drbear95/gonotter-server/utils"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
)

func init() {
	log.Println("Mongo connection initialization")

	c := utils.Config{}
	c.GetConfig()

	url := "mongodb://localhost:27017/?readPreference=primary&appname=gonotterdb"
	_ = mgm.SetDefaultConfig(nil, "gonotterdb", options.Client().ApplyURI(url))

	_, err := mgm.Coll(&model.Note{}).Indexes().CreateOne(mgm.Ctx(), mongo.IndexModel{
		Keys: bson.M{"title": bsonx.String("text"), "content": bsonx.String("text")},
	})

	if err != nil {
		log.Println(err)
	}
}
