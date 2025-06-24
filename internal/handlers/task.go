package handlers

import (
	"fmt"
	"task_manager/internal/models"
	"task_manager/internal/repositories"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

func GetTasks(collection *mongo.Collection) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		user := c.Locals("user").(*models.User)
		r := repositories.NewTaskRepository(collection)
		tasks, err := r.GetTasks(user, ctx)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("Error: %v", err)})
		}
		return c.Status(200).JSON(fiber.Map{"tasks": tasks})
	}
}

func CreateTask(collection *mongo.Collection) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		user := c.Locals("user").(*models.User)
		task := new(models.Task)
		if err := c.BodyParser(task); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Invalid request body"})
		}

		validate := validator.New()
		if err := validate.Struct(task); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": err.Error()})
		}
		task.UserID = user.ID
		r := repositories.NewTaskRepository(collection)
		if _, err := r.CreateTask(task, ctx); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("Error: %v", err)})
		}
		return c.Status(200).JSON(fiber.Map{"message": "Task created successfully"})
	}
}

func EditTask(collection *mongo.Collection) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		user := c.Locals("user").(*models.User)
		task := new(models.Task)
		if err := c.BodyParser(task); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Invalid request body"})
		}

		validate := validator.New()
		if err := validate.Struct(task); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": err.Error()})
		}

		task.UserID = user.ID
		r := repositories.NewTaskRepository(collection)
		if _, err := r.UpdateTask(task, ctx); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("Error: %v", err)})
		}
		return c.Status(200).JSON(fiber.Map{"message": "Task edited successfully"})
	}
}

func DeleteTask(collection *mongo.Collection) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		userID := (c.Locals("user").(*models.User)).ID

		taskID, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Invalid task ID"})
		}
		r := repositories.NewTaskRepository(collection)
		if _, err := r.DeleteTask(userID, taskID, ctx); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("Error: %v", err)})
		}
		return c.Status(200).JSON(fiber.Map{"message": "Task deleted successfully"})
	}
}
