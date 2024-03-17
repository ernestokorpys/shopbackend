package controller

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/ernestokorpys/shopbackend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FilterProducts(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse Body")
		return err
	}

	client := c.Locals("mongoClient").(*mongo.Client)
	collection := client.Database("onlineshop").Collection("products")

	filter := bson.M{}

	// Añadir filtros según los parámetros recibidos
	if minCostStr, ok := data["minCost"].(string); ok {
		minCost, err := strconv.ParseInt(minCostStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid minCost",
			})
		}
		filter["cost"] = bson.M{"$gte": minCost}
	}

	if maxCostStr, ok := data["maxCost"].(string); ok {
		maxCost, err := strconv.ParseInt(maxCostStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid maxCost",
			})
		}
		filter["cost"] = bson.M{"$lte": maxCost}
	}

	if keyword, ok := data["keyword"].(string); ok {
		filter["productName"] = bson.M{"$regex": keyword, "$options": "i"}
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page number",
		})
	}

	limit := int64(4)
	offset := int64((page - 1) * 4)

	// Realizar la consulta con el filtro para obtener los productos paginados
	cursor, err := collection.Find(context.TODO(), filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &offset,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error finding products",
		})
	}
	defer cursor.Close(context.TODO())

	var products []models.Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error decoding products",
		})
	}

	// Obtener la cantidad total de productos que cumplen con el filtro
	total, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error counting total products",
		})
	}

	// Calcular el número total de páginas
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Devolver la respuesta JSON con los productos, la información de paginación y la cantidad total de productos
	return c.JSON(fiber.Map{
		"data": products,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": totalPages,
		},
	})
}
