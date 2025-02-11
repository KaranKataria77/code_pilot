package controllers

import (
	"code_pilot/internal/config"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUser(c *gin.Context) {
	var user struct {
		ID    string `bson:"_id"`
		Email string `bson:"email"`
		Name  string `bson:"name"`
	}
	collection := config.GetMongoCollection("users")
	ID := c.Param("id")
	bsonId, bson_err := primitive.ObjectIDFromHex(ID)
	if bson_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while converting string to bson " + bson_err.Error()})
	}
	err := collection.FindOne(context.TODO(), bson.M{"_id": bsonId}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching user" + err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "User profile ", "user": user})
}
