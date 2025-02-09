package controllers

import (
	"code_pilot/internal/config"
	"code_pilot/internal/models"
	"code_pilot/internal/utils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func SignUp(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// check if user already exists
	collection := config.GetMongoCollection("users")
	var existingUser models.User
	existing_user_err := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if existing_user_err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	// hash the password
	hashedPassword, hashing_err := utils.HashPassword(user.Password)
	if hashing_err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	// save user to database
	_, insertion_err := collection.InsertOne(context.TODO(), user)
	if insertion_err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting user to database"})
		return
	}

	// get jwt token
	token, token_err := utils.GenerateJWTToken(user.Email)
	if token_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while generating token "})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User insertion successfully", "token": token})
}

func Login(c *gin.Context) {
	type LoginUser struct {
		Email    string `bson:"email" binding:"required,email"`
		Password string `bson:"password" binding:"required"`
	}

	var user LoginUser

	err := c.ShouldBindJSON(&user)
	passwordString := user.Password
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide valid Email or password"})
		return
	}

	// find user
	collection := config.GetMongoCollection("users")
	user_err := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&user)
	if user_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// validate password
	is_password_same := utils.ComparePassword(user.Password, passwordString)
	if !is_password_same {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not valid"})
		return
	}

	// generate token
	token, token_err := utils.GenerateJWTToken(user.Email)
	if token_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while generating token "})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Login successfully", "token": token})

}
