package main

import (
	"os"

	"github.com/ernestokorpys/shopbackend/database"
	"github.com/ernestokorpys/shopbackend/middleware"
	"github.com/ernestokorpys/shopbackend/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	client := database.ConnectDB() //cliente de mongo.DB
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error")
	// }
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
