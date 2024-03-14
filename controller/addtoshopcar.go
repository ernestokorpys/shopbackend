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

type RequestData struct {
	ProductIDs []string `json:"productIDs"`
}

func AddToShopCar(c *fiber.Ctx) error {
	// Obtener el ID del usuario autenticado

	var data RequestData // Parsear los datos de la solicitud
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).SendString("Error al analizar el cuerpo de la solicitud")
	}

	client := c.Locals("mongoClient").(*mongo.Client)
	// // Obtener la colección de usuarios
	collection := client.Database("onlineshop").Collection("shopUsers")
	// User ID search database ---------------------------------------------
	cookie := c.Cookies("jwt")
	userID, _ := util.Parsejwt(cookie)

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).SendString("ID de usuario inválido")
	}

	// Buscar el usuario por su ID
	filter := bson.M{"_id": objectID}
	var user models.User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return c.Status(404).SendString("Usuario no encontrado")
	}

	var productIDs []primitive.ObjectID
	for _, id := range data.ProductIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(400).SendString("ID de producto inválido")
		}
		productIDs = append(productIDs, objectID)
	}
	// Obtener los productos correspondientes a los IDs
	productsCollection := client.Database("onlineshop").Collection("products")
	cursor, err := productsCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": productIDs}})
	if err != nil {
		return c.Status(500).SendString("Error al obtener los productos")
	}
	defer cursor.Close(context.Background())

	// Iterar sobre los productos y agregarlos al carrito
	for cursor.Next(context.Background()) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			return c.Status(500).SendString("Error al decodificar los productos")
		}
		user.ShopCar.Products = append(user.ShopCar.Products, product)
	}

	// Actualizar el usuario en la base de datos
	update := bson.M{"$set": bson.M{"shopCar": user.ShopCar}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(500).SendString("Error al actualizar el usuario")
	}

	// Convertir el ID del usuario a ObjectID

	return c.JSON(fiber.Map{
		"user":    userID,
		"data":    data,
		"message": "Productos agregados al carrito con éxito",
	})
}
