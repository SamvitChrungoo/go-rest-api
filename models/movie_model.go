package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//Movie model for DB
type Movie struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title         string             `json:"title,omitempty" bson:"title,omitempty"`
	Year          int                `json:"year,omitempty" bson:"year,omitempty"`
	OriginalTitle string             `json:"originalTitle,omitempty" bson:"originalTitle,omitempty"`
	Storyline     string             `json:"storyline,omitempty" bson:"storyline,omitempty"`
	ImdbRating    float64            `json:"imdbRating,omitempty" bson:"imdbRating,omitempty"`
	PosterURL     string             `json:"posterurl,omitempty" bson:"posterurl,omitempty"`
}

//ErrorResponse -> for displayin errors
type ErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
