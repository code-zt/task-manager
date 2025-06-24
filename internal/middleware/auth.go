package middleware

import (
	"context"
	"strings"
	"task_manager/internal/repositories"
	"task_manager/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthMiddleware(collection *mongo.Collection) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		refreshToken := c.Cookies("refreshToken")
		if after, ok := strings.CutPrefix(refreshToken, "Bearer "); ok {
			refreshToken = after
		}
		if refreshToken == "" {
			return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
		}

		accessToken := c.Cookies("accessToken")
		if after, ok := strings.CutPrefix(accessToken, "Bearer "); ok {
			accessToken = after
		}
		if accessToken == "" {
			return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
		}

		user, err := utils.ValidateAccessToken(accessToken)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
		}

		_, err = utils.ValidateRefreshToken(refreshToken)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
		}

		r := repositories.NewUserRepository(collection)
		existingUser, err := r.FindUserByEmail(user.Email, ctx)
		if err != nil || existingUser == nil {
			return c.Status(404).JSON(fiber.Map{"message": "User  not found"})
		}

		c.Locals("user", existingUser)

		return c.Next()
	}
}
