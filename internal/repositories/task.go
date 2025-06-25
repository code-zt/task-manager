package repositories

import (
	"context"
	"fmt"
	"task_manager/internal/config"
	"task_manager/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var cfg = config.LoadConfig()

func NewTaskRepository(db *mongo.Collection) *TaskRepository {
	return &TaskRepository{db: db}
}

type TaskRepository struct {
	db *mongo.Collection
}

func (t *TaskRepository) CreateTask(task *models.Task, ctx context.Context) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.ContextTimeout)
	defer cancel()
	task.CreatedAt = time.Now()
	result, err := t.db.InsertOne(ctx, task)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *TaskRepository) GetTasks(user *models.User, ctx context.Context) ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.ContextTimeout)
	defer cancel()
	var tasks []models.Task
	cursor, err := t.db.Find(ctx, bson.M{"user_id": user.ID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx) // Закрываем курсор после использования
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *TaskRepository) GetTask(taskID primitive.ObjectID, user *models.User, ctx context.Context) (*models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.ContextTimeout)
	defer cancel()
	var task models.Task
	err := t.db.FindOne(ctx, bson.M{"_id": taskID, "user_id": user.ID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return &task, nil
}

func (t *TaskRepository) UpdateTask(task *models.Task, ctx context.Context) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.ContextTimeout)
	defer cancel()
	result, err := t.db.UpdateOne(ctx, bson.M{"_id": task.ID, "user_id": task.UserID}, bson.M{"$set": task})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *TaskRepository) DeleteTask(userID, taskID primitive.ObjectID, ctx context.Context) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.ContextTimeout)
	defer cancel()
	result, err := t.db.DeleteOne(ctx, bson.M{"_id": taskID, "user_id": userID})
	if err != nil {
		return nil, err
	}
	return result, nil
}
