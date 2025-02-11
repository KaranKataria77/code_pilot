package controllers

import (
	"code_pilot/internal/config"
	"code_pilot/internal/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProject(c *gin.Context) {
	ID := c.Param("userID")
	userBsonId, user_bson_err := primitive.ObjectIDFromHex(ID)
	if user_bson_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while converting string to bson ID " + user_bson_err.Error()})
		return
	}
	var fetch_project struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	created_at := time.Now()
	updated_at := time.Now()

	err := c.ShouldBindJSON(&fetch_project)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid fields " + err.Error()})
		return
	}
	var project models.Project
	project.Name = fetch_project.Name
	project.Description = fetch_project.Description
	project.OwnerID = userBsonId
	project.CreatedAt = created_at
	project.UpdatedAt = updated_at
	collection := config.GetMongoCollection("projects")
	_, project_err := collection.InsertOne(context.TODO(), project)
	if project_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while inserting project " + project_err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully", "project": project})
}
