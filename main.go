package main

import (
	"os"

	"github.com/ernestokorpys/shopbackend/database"
	"github.com/ernestokorpys/shopbackend/middleware"
	"github.com/ernestokorpys/shopbackend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // Importa el middleware CORS de Fiber
)

func main() {
	client := database.ConnectDB() //cliente de mongo.DB
	/*
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	*/
	port := os.Getenv("PORT")
	app := fiber.New()
	// Middleware de recuperación de errores
	// Middleware de manejo de errores
	app.Use(middleware.ErrorHandler())

	// Middleware para establecer el cliente de MongoDB en el contexto
	app.Use(middleware.MongoDBClientSetter(client))

	// Configuración de CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173", // Permite solicitudes desde este origen
		AllowMethods:     "GET,POST,PUT,DELETE",   // Permite estos métodos HTTP
		AllowCredentials: true,                    // Permite credenciales (cookies)

	}))

	routes.Setup(app)      //maneja las solicitudes entrantes
	app.Listen(":" + port) //Escucha las solicitudes del puerto constantemente
}
