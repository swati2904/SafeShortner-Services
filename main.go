package main

import (
	"SafeShotner-Services/configs"
	"SafeShotner-Services/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	//initialize a fiber application
	app := fiber.New()

	// Enable CORS for all routes
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	//Run database
	configs.ConnectDB()

	//routes
	routes.RegisterURLRoutes(app)

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.JSON(&fiber.Map{"data": "Hello from fiber & mongoDB"})
	// })
	app.Listen(":5000")
}
