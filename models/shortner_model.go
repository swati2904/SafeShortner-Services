package models

type Shortner struct {
	ID       int    `json:"id" bson:"id" validate:"required"`
	LongURL  string `json:"longURL" bson:"longURL" validate:"required"`
	ShortURL string `json:"shortURL" bson:"shortURL" validate:"required"`
}
