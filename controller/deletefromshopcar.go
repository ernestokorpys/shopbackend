package controller

import (
	"context"
	"strconv"

	"github.com/ernestokorpys/shopbackend/models"
	"github.com/ernestokorpys/shopbackend/util"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RemoveFromShopCar(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	userID, _ := util.Parsejwt(cookie)
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).SendString("ID de usuario inválido")
	}

	client := c.Locals("mongoClient").(*mongo.Client)
	collection := client.Database("onlineshop").Collection("shopUsers")

	// Obtener el índice del producto en la lista
	indexToDelete, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Índice de producto inválido")
	}

	// Buscar el usuario por su ID
	filter := bson.M{"_id": objectID}
	var user models.User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return c.Status(404).SendString("Usuario no encontrado")
	}

	// Verificar si el índice está dentro del rango de la lista de productos
	if indexToDelete < 0 || indexToDelete >= len(user.ShopCar.Products) {
		return c.Status(400).SendString("Índice de producto fuera de rango")
	}

	// Eliminar el producto del carrito de compras del usuario en la posición especificada
	user.ShopCar.Products = append(user.ShopCar.Products[:indexToDelete], user.ShopCar.Products[indexToDelete+1:]...)

	// Actualizar el usuario en la base de datos
	update := bson.M{"$set": bson.M{"shopCar": user.ShopCar}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(500).SendString("Error al actualizar el usuario")
	}

	return c.JSON(fiber.Map{
		"message": "Producto eliminado del carrito con éxito",
	})
}
