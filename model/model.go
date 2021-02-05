package model

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccessTokenDetails struct {
	mgm.DefaultModel `bson:",inline"`
	UserID           string `json:"user_id" bson:"user_id"`
	Uuid             string `json:"uuid" bson:"uuid"`
	ExpiresAt        int64  `json:"expires_at" bson:"expires_at"`
}

type RefreshTokenDetails struct {
	mgm.DefaultModel `bson:",inline"`
	UserID           string `json:"user_id" bson:"user_id"`
	Uuid             string `json:"uuid" bson:"uuid"`
	ExpiresAt        int64  `json:"expires_at" bson:"expires_at"`
}

type Note struct {
	mgm.DefaultModel `bson:",inline"`
	Title            string             `json:"title" bson:"title"`
	Content          string             `json:"content" bson:"content"`
	AuthorID         primitive.ObjectID `json:"author_id" bson:"author_id"`
}

type User struct {
	mgm.DefaultModel `json:"-" bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Password         string `json:"password" bson:"password"`
	Email            string `json:"email" bson:"email"`
}

func NewNote(title string, content string, authID primitive.ObjectID) *Note {
	return &Note{
		Title:    title,
		Content:  content,
		AuthorID: authID,
	}
}

func NewUser(name string, password string, email string) *User {
	return &User{
		Name:     name,
		Password: password,
		Email: email,
	}
}

func NewRefreshTokenDetails(userId string, uuid string, expiresAt int64) *RefreshTokenDetails {
	return &RefreshTokenDetails{
		UserID:    userId,
		Uuid:      uuid,
		ExpiresAt: expiresAt,
	}
}

func NewAccessTokenDetails(userId string, uuid string, expiresAt int64) *AccessTokenDetails {
	return &AccessTokenDetails{
		UserID:    userId,
		Uuid:      uuid,
		ExpiresAt: expiresAt,
	}
}
