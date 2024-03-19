package controller

import (
	"context"
	"math"
	"strconv"

	"github.com/ernestokorpys/shopbackend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FilterProducts(c *fiber.Ctx) error {
	client := c.Locals("mongoClient").(*mongo.Client)
	collection := client.Database("onlineshop").Collection("products")

	filter := bson.M{}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page number",
		})
	}

	minCost, err := strconv.Atoi(c.Query("minCost", "0"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid min cost number",
		})
	}
	filter["cost"] = bson.M{"$gte": minCost}

	maxCostStr := c.Query("maxCost", "100000000000") // Valor predeterminado si no se proporciona
	maxCost, err := strconv.ParseInt(maxCostStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid max cost number",
		})
	}
	// Combina los filtros minCost y maxCost para el costo
	filter["cost"].(bson.M)["$lte"] = maxCost

	keyword := c.Query("keyword", "")
	if keyword != "" {
		// Añade el filtro para el nombre del producto utilizando una expresión regular
		filter["productName"] = bson.M{"$regex": keyword, "$options": "i"}
	}

	limit := int64(8)
	offset := int64((page - 1) * 4)

	// Realiza la consulta con el filtro para obtener los productos paginados
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
