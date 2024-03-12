package routes

import (
	"github.com/ernestokorpys/shopbackend/controller"
	"github.com/ernestokorpys/shopbackend/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controller.RegisterUser)
	app.Post("/api/login", controller.Login)
	app.Use(middleware.IsAuthenticate)
	app.Post("/api/addproduct", controller.AddProduct)
	app.Get("/api/products", controller.GetProducts)
	app.Post("/api/logout", controller.Logout) // Nueva ruta para cerrar sesión
}
