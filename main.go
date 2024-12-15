package main

import (
	"SafeShotner-Services/configs"
	"SafeShotner-Services/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	//initialize a fiber application
	app := fiber.New()

	//Run database
	configs.ConnectDB()

	//routes
	routes.RegisterURLRoutes(app)

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.JSON(&fiber.Map{"data": "Hello from fiber & mongoDB"})
	// })
	app.Listen(":6000")
}
