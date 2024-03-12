package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	// Remover la cookie del token JWT del cliente
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie) // Agregar la cookie a la
	// Puedes enviar una respuesta indicando que el logout fue exitoso
	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}
