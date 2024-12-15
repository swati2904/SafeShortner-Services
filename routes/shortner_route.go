package routes

import (
	"github.com/gofiber/fiber/v2"

	"SafeShotner-Services/controllers"
)

func RegisterURLRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/shorten", controllers.CreateShortURL)
}
