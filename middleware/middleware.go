package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Recuperación de un panic
				log.Printf("Recovered from panic: %v", r)
				c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error. bad JSON")
			}
		}()
		return c.Next()
	}
}

// Middleware para establecer el cliente de MongoDB en el contexto
func MongoDBClientSetter(client *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Establecer el cliente de MongoDB en el contexto
		c.Locals("mongoClient", client)
		return c.Next()
	}
}
