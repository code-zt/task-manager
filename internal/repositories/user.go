package repositories

import (
	"context"
	"fmt"
	"task_manager/internal/models"
	"task_manager/internal/utils"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Collection) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) CreateUser(user *models.User, ctx context.Context) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.ContextTimeout)
	defer cancel()
	user.CreatedAt = time.Now()
	result, err := u.db.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *UserRepository) FindUserByEmail(email string, ctx context.Context) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.ContextTimeout)
	defer cancel()
	var user models.User
	err := u.db.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) Auth(email, password string, ctx context.Context) (*models.User, error) {
	var user models.User

	err := u.db.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &user, nil
}

func (u *UserRepository) FindUserByID(id string, ctx context.Context) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.ContextTimeout)
	defer cancel()
	var user models.User
	err := u.db.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
