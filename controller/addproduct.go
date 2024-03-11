package controller

import (
	"context"
	"fmt"

	"github.com/ernestokorpys/shopbackend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddProduct(c *fiber.Ctx) error {
	// Parsear los datos JSON de la solicitud
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Parsing request body failed",
		})
	}

	// Verificar si hay campos adicionales en el JSON
	expectedKeys := map[string]bool{
		"productName": true,
		"cost":        true,
		"picture":     true,
	}
	for key := range data {
		if !expectedKeys[key] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Unexpected field '%s' in request body", key),
			})
		}
	}

	// Crear un nuevo usuario con los datos recibidos
	user := models.Product{
		ProductName: data["productName"].(string),
		Cost:        data["cost"].(string),
		Picture:     data["picture"].(string),
	}

	// Obtener el cliente de MongoDB de la variable de contexto
	client := c.Locals("mongoClient").(*mongo.Client)

	// Obtener la colección de usuarios
	collection := client.Database("onlineshop").Collection("products")

	// Insertar el nuevo usuario en la base de datos
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert user into database",
		})
	}

	// Retornar una respuesta de éxito
	c.Status(200)
	return c.JSON(fiber.Map{
		"user":    user,
		"message": "Product added successfully.",
	})
}
