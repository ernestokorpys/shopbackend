package controller

import (
	"context"
	"fmt"

	"github.com/ernestokorpys/shopbackend/utils"

	"github.com/ernestokorpys/shopbackend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(c *fiber.Ctx) error {
	// Parsear los datos JSON de la solicitud
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Parsing request body failed",
		})
	}

	// Verificar si hay campos adicionales en el JSON
	var stringData map[string]string
	for key, value := range data {
		stringData[key] = fmt.Sprintf("%v", value)
	}

	expectedFields := []string{"userName", "email", "password"}
	if err := utils.CheckExpectedFields(stringData, expectedFields); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Obtener el cliente de MongoDB de la variable de contexto
	client := c.Locals("mongoClient").(*mongo.Client)

	// Obtener la colección de usuarios
	collection := client.Database("onlineshop").Collection("shopUsers")

	// Verificar si el correo electrónico ya existe en la base de datos
	existingUser := models.User{}
	err := collection.FindOne(context.Background(), bson.M{"email": data["email"].(string)}).Decode(&existingUser)
	if err == nil {
		// El correo electrónico ya está en uso, retornar un error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Email already exists in the database",
			"warning": "Failed try to acces acound send to" + existingUser.Email,
			"user":    existingUser,
		})
	}

	// Crear un nuevo usuario con los datos recibidos
	newUser := models.User{
		UserName: data["userName"].(string),
		Email:    data["email"].(string),
	}

	newUser.EncryptedPassword(data["password"].(string))

	// Insertar el nuevo usuario en la base de datos
	_, err = collection.InsertOne(context.Background(), newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert user into database",
		})
	}

	// Retornar una respuesta de éxito
	return c.JSON(fiber.Map{
		"user":    newUser,
		"message": "Account created successfully.",
	})
}
