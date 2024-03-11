package routes

import (
	"github.com/ernestokorpys/shopbackend/controller"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controller.RegisterUser)
	app.Post("/api/addproduct", controller.AddProduct)
	app.Get("/api/products", controller.GetProducts)
	app.Post("/api/login", controller.Login)

	// app.Use(middleware.IsAuthenticate)
	// app.Post("/api/post", controller.CreatePost)
	// app.Get("/api/allpost", controller.AllPost)
	// app.Get("/api/allpost/:id", controller.DetailPost)
	// app.Put("/api/updatepost/:id", controller.UpdatePost)
	// app.Get("/api/uniquepost", controller.UniquePost)
	// app.Delete("/api/deletepost/:id", controller.DeletePost)
	// app.Post("/api/upload-image", controller.UploadImage)
	// app.Static("/api/uploads", "./uploads")

}
