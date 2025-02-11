package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID            string               `bson:"_id"`
	Name          string               `bson:"name" binding:"required"`
	OwnerID       primitive.ObjectID   `bson:"owner" binding:"required"`
	Collaborators []primitive.ObjectID `bson:"collaborators,omitempty"`
	Description   string               `bson:"description"`
	CreatedAt     time.Time            `bson:"created_at" binding:"required"`
	UpdatedAt     time.Time            `bson:"updated_at" binding:"required"`
}
