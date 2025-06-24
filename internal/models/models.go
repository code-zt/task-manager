package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title       string             `json:"title" bson:"title" validate:"required"`
	Description string             `json:"description" bson:"description"`
	Status      string             `json:"status" bson:"status" validate:"required,oneof=pending in_progress completed"`
	Priority    string             `json:"priority" bson:"priority" validate:"required,oneof=low medium high"`
	DueDate     *time.Time         `json:"due_date" bson:"due_date"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username" validate:"required,min=3"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"-" bson:"password" validate:"required,min=6"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
