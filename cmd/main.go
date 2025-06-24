package main

import (
	"context"
	"log"
	"task_manager/internal/config"
	"task_manager/internal/database"
	"task_manager/internal/handlers"
	"task_manager/internal/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var cfg *config.Config
var client, err = database.Connect()

var userCollection *mongo.Collection = client.Database.Collection("users")
var taskCollection *mongo.Collection = client.Database.Collection("tasks")

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer func() {

		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting from database: %v", err)
		} else {
			log.Println("Disconnected from database")
		}
	}()

	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	app := fiber.New()

	app.Use(middleware.RequestLogger(), middleware.ErrorLogger())
	api := app.Group("/api")
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"message": "pong"})
	})
	api.Post("/Register", handlers.Register(userCollection))
	api.Post("/Login", handlers.Login(userCollection))
	api.Post("/Logout", handlers.Logout)

	api.Use(middleware.AuthMiddleware(userCollection))
	task := api.Group("/task")
	task.Post("/create", handlers.CreateTask(taskCollection))
	task.Get("/get", handlers.GetTasks(taskCollection))
	task.Put("/edit", handlers.EditTask(taskCollection))
	task.Delete("/delete", handlers.DeleteTask(taskCollection))

	app.Listen(":3000")
}
