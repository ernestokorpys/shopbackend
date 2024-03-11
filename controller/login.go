package controller

import (
	"context"
	"fmt"

	"github.com/ernestokorpys/shopbackend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse Body")
	}
	// Obtener el cliente de MongoDB de la variable de contexto
	client := c.Locals("mongoClient").(*mongo.Client)
	// Obtener la colección de usuarios
	collection := client.Database("onlineshop").Collection("shopUsers")

	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": data["email"]}).Decode(&user)
	if err != nil {
		// El correo electrónico no se encontró en la base de datos
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Correo no encontrado en la base de datos",
		})
	}
	if err := user.ComparePassword(data["password"]); err != nil {
		// La contraseña no coincide
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Contraseña incorrecta",
		})
	}

	return c.JSON(fiber.Map{
		"massage": "Succesful login",
		"user":    user,
	})

}
