package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GlicemiaInput struct {
	Glucose float64 `json:"glucose" binding:"required"`
	Time    string  `json:"time" binding:"required"`
}

type GlicemiaRecord struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Person    string             `bson:"person" json:"person"`
	Glucose   float64            `bson:"glucose" json:"glucose"`
	Date      string             `bson:"date" json:"date"`
	Time      string             `bson:"time" json:"time"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
