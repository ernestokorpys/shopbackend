package controller

import (
	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	// Eliminar la cookie de sesión o invalidar el token de acceso
	c.ClearCookie("jwt")
	// Opcionalmente, puedes redirigir al usuario a una página de inicio de sesión u otra página después del log out
	return c.SendString("Logged out successfully")
}
