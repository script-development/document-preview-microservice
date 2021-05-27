package main

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func getKey() string {
	return strings.TrimSpace(os.Getenv("KEY"))
}

func requireAuth() fiber.Handler {
	key := getKey()
	return func(c *fiber.Ctx) error {
		auth := strings.TrimSpace(c.Get("Authorization"))
		if len(auth) == 0 || !strings.HasPrefix(auth, "Bearer ") {
			c.Set("WWW-Authenticate", "Bearer")
			return c.SendStatus(401)
		}

		auth = strings.TrimPrefix(auth, "Bearer ")
		if auth != key {
			c.Set("WWW-Authenticate", "Bearer")
			return c.SendStatus(401)
		}

		return c.Next()
	}
}
