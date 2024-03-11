package controller

import (
	"context"
	"math"
	"strconv"

	"github.com/ernestokorpys/shopbackend/models"
	"github.com/ernestokorpys/shopbackend/util"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetProducts(c *fiber.Ctx) error {
	// Obtener el token de autorización del encabezado de la solicitud
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No se proporcionó el token de autorización",
		})
	}

	// Verificar el token de autorización
	_, err := util.VerifyJwt(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token de autorización inválido",
		})
	}

	page, err := strconv.Atoi(c.Query("page", "1")) //se fija el parametro page en la url si no tine asume que es 1
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page number",
		})
	}

	limit := 2
	offset := (page - 1) * limit

	client := c.Locals("mongoClient").(*mongo.Client)

	// Obtener la colección de productos
	collection := client.Database("onlineshop").Collection("products")

	// Obtener la cantidad total de productos
	total, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count products",
		})
	}

	// Consultar los productos paginados debido a que cada limite elementos, debera esperar y reorganizarlo en grupos
	cursor, err := collection.Find(context.Background(), bson.M{}, options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}
	defer cursor.Close(context.Background())

	var getproducts []models.Product
	if err := cursor.All(context.Background(), &getproducts); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode products",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if page > totalPages {
		return c.JSON(fiber.Map{
			"warning": "Page number exceeds total pages",
		})
	}

	return c.JSON(fiber.Map{
		"data": getproducts,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": math.Ceil(float64(total) / float64(limit)),
		},
	})
}
