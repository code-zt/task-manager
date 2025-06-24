package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		req := c.Request()
		method := string(req.Header.Method())
		path := string(req.URI().Path())
		ip := c.IP()

		err := c.Next()

		status := c.Response().StatusCode()
		latency := time.Since(start)
		errorMsg := ""
		if err != nil {
			errorMsg = err.Error()
		}

		log.Info().
			Str("method", method).
			Str("path", path).
			Str("ip", ip).
			Int("status", status).
			Dur("latency", latency).
			Str("error", errorMsg).
			Msg("handled request")

		return err
	}
}

func ErrorLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			log.Error().
				Str("path", c.Path()).
				Int("status", c.Response().StatusCode()).
				Str("error", err.Error()).
				Msg("request failed")
		}
		return err
	}
}
