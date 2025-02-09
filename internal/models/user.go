package models

type User struct {
	ID       string `bson:"_id,omitempty"`
	Name     string `bson:"name" binding:"required"`
	Email    string `bson:"email" binding:"required,email"`
	Password string `bson:"password" binding:"required"`
}
