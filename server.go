package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/go-fiber-bookstore-auth-api/routes"
)

func main() {
	app := fiber.New()
	routes.AuthRoutes(app)
	app.Listen(":9020")
}
