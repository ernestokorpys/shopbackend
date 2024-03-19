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
		AllowOrigins:     "http://localhost:5173,https://frontendreact-neon.vercel.app,https://frontendreact-92ossl6rw-ernestokorpys.vercel.app,https://frontendreact-ernestokorpys.vercel.app,https://frontendreact-git-master-ernestokorpys.vercel.app",
		AllowMethods:     "GET,POST,PUT,DELETE", // Permite estos métodos HTTP
		AllowCredentials: true,                  // Permite credenciales (cookies)
		AllowHeaders:     "Content-Type",        // Permitir los encabezados necesarios
	}))

	routes.Setup(app)      //maneja las solicitudes entrantes
	app.Listen(":" + port) //Escucha las solicitudes del puerto constantemente
}
