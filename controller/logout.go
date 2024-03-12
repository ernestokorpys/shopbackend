package controller

import "github.com/gofiber/fiber/v2"

func Logout(c *fiber.Ctx) error {
	// Remover la cookie del token JWT del cliente
	c.ClearCookie("jwt")

	// Puedes enviar una respuesta indicando que el logout fue exitoso
	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}
