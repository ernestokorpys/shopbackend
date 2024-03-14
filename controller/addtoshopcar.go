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
	cookie := c.Cookies("jwt")
	userID, _ := util.Parsejwt(cookie)
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).SendString("ID de usuario inválido")
	}

	// Parsear los datos de la solicitud
	var data struct {
		ProductID string `json:"productID"`
	}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).SendString("Error al analizar el cuerpo de la solicitud")
	}

	client := c.Locals("mongoClient").(*mongo.Client)
	collection := client.Database("onlineshop").Collection("shopUsers")

	// Buscar el usuario por su ID
	filter := bson.M{"_id": objectID}
	var user models.User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return c.Status(404).SendString("Usuario no encontrado")
	}

	// Convertir el ID del producto a ObjectID
	productID, err := primitive.ObjectIDFromHex(data.ProductID)
	if err != nil {
		return c.Status(400).SendString("ID de producto inválido")
	}

	// Obtener el producto correspondiente al ID
	productsCollection := client.Database("onlineshop").Collection("products")
	var product models.Product
	err = productsCollection.FindOne(context.Background(), bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		return c.Status(404).SendString("Producto no encontrado")
	}

	// Agregar el producto al carrito del usuario
	user.ShopCar.Products = append(user.ShopCar.Products, product)

	// Actualizar el usuario en la base de datos
	update := bson.M{"$set": bson.M{"shopCar": user.ShopCar}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(500).SendString("Error al actualizar el usuario")
	}

	return c.JSON(fiber.Map{
		"message": "Producto agregado al carrito con éxito",
	})
}
