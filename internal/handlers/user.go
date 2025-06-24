package handlers

import (
	"fmt"
	"task_manager/internal/config"
	"task_manager/internal/models"
	"task_manager/internal/repositories"
	"task_manager/internal/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var cfg, _ = config.LoadConfig()

var (
	accessTokenLifetime  = 15 * time.Minute
	refreshTokenLifetime = 30 * 24 * time.Hour
)

func Register(collection *mongo.Collection) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": fmt.Sprintf("Error: %v", err)})
		}

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": err.Error()})
		}
		r := repositories.NewUserRepository(collection)
		existingUser, err := r.FindUserByEmail(user.Email, ctx)
		if err != mongo.ErrNoDocuments && err != nil {
			return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("Error: %v", err)})
		}
		if existingUser != nil {
			return c.Status(400).JSON(fiber.Map{"message": "User  already exists"})
		}

		user.Password, err = utils.HashPassword(user.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Cannot hash password"})
		}

		result, err := r.CreateUser(&user, ctx)
		if err != nil || result == nil {
			return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("Error: %v", err)})
		}

		accessToken, err := utils.CreateAccessToken(&user)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("Error: %v", err)})
		}

		refreshToken, err := utils.CreateRefreshToken(&user)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("Error: %v", err)})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "refreshToken",
			Path:     "/",
			Value:    refreshToken,
			Secure:   cfg.UseHttps,
			HTTPOnly: true,
			Expires:  time.Now().Add(refreshTokenLifetime),
		})
		c.Cookie(&fiber.Cookie{
			Name:     "accessToken",
			Path:     "/",
			Value:    accessToken,
			Secure:   cfg.UseHttps,
			HTTPOnly: true,
			Expires:  time.Now().Add(accessTokenLifetime),
		})

		return c.Status(201).JSON(fiber.Map{"message": "User  created", "id": result.InsertedID, "accessToken": accessToken, "refreshToken": refreshToken})
	}
}

func Login(collection *mongo.Collection) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		r := repositories.NewUserRepository(collection)

		var user models.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": fmt.Sprintf("Error: %v", err)})
		}

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": err.Error()})
		}

		authUser, err := r.Auth(user.Email, user.Password, ctx)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"message": "Invalid credentials"})
		}

		accessToken, err := utils.CreateAccessToken(authUser)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("Error: %v", err)})
		}
		refreshToken, err := utils.CreateRefreshToken(authUser)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": fmt.Sprintf("Error: %v", err)})
		}
		c.Cookie(&fiber.Cookie{
			Name:     "refreshToken",
			Path:     "/",
			Value:    refreshToken,
			Secure:   cfg.UseHttps,
			HTTPOnly: true,
			Expires:  time.Now().Add(refreshTokenLifetime),
		})
		c.Cookie(&fiber.Cookie{
			Name:     "accessToken",
			Path:     "/",
			Value:    accessToken,
			Secure:   cfg.UseHttps,
			HTTPOnly: true,
			Expires:  time.Now().Add(accessTokenLifetime),
		})
		return c.Status(200).JSON(fiber.Map{"message": "Login successful", "refreshToken": refreshToken, "accessToken": accessToken})
	}
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("refreshToken", "accessToken")
	return c.Status(200).JSON(fiber.Map{"message": "Logout successful"})
}
