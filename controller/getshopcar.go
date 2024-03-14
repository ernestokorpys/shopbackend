package controller

import (
	"context"

	"github.com/ernestokorpys/shopbackend/models"
	"github.com/ernestokorpys/shopbackend/util"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetShopCar(c *fiber.Ctx) error {
	client := c.Locals("mongoClient").(*mongo.Client)
	// // Obtener la colección de usuarios
	collection := client.Database("onlineshop").Collection("shopUsers")

	cookie := c.Cookies("jwt")
	userID, _ := util.Parsejwt(cookie)

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).SendString("ID de usuario inválido")
	}

	filter := bson.M{"_id": objectID}
	var user models.User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return c.Status(404).SendString("Usuario no encontrado")
	}

	return c.JSON(fiber.Map{
		"data":    user.ShopCar.Products,
		"test":    user.UserName,
		"test2":   user.Email,
		"Shopcar": user.ShopCar,
	})
}
