package main

import (
	"log"
	"os"

	"github.com/ernestokorpys/shopbackend/database"
	"github.com/ernestokorpys/shopbackend/middleware"
	"github.com/ernestokorpys/shopbackend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	client := database.ConnectDB() //cliente de mongo.DB
	// Colección y contexto
	// collection := client.Database("onlineshop").Collection("shopUsers")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error")
	}
	port := os.Getenv("PORT")
	app := fiber.New()
	// Middleware de recuperación de errores
	// Middleware de manejo de errores
	app.Use(middleware.ErrorHandler())

	// Middleware para establecer el cliente de MongoDB en el contexto
	app.Use(middleware.MongoDBClientSetter(client))

	routes.Setup(app)      //maneja las solicitudes entrantes
	app.Listen(":" + port) //Escucha las solicitudes del puerto constantemente
}
