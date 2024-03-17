package controller

import (
	"context"
	"fmt"
	"strconv"

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

	// Convertir el campo "cost" del JSON a un entero
	cost, err := strconv.ParseInt(data["cost"].(string), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid value for 'cost'",
		})
	}

	// Crear un nuevo producto con los datos recibidos
	product := models.Product{
		ProductName: data["productName"].(string),
		Cost:        cost,
		Picture:     data["picture"].(string),
	}

	// Obtener el cliente de MongoDB de la variable de contexto
	client := c.Locals("mongoClient").(*mongo.Client)

	// Obtener la colección de productos
	collection := client.Database("onlineshop").Collection("products")

	// Insertar el nuevo producto en la base de datos
	_, err = collection.InsertOne(context.Background(), product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert product into database",
		})
	}

	// Retornar una respuesta de éxito
	return c.JSON(fiber.Map{
		"product": product,
		"message": "Product added successfully.",
	})
}
