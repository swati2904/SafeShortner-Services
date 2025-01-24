package routes

import (
	"github.com/gofiber/fiber/v2"

	"SafeShotner-Services/controllers"
)

func RegisterURLRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/shorten", controllers.CreateShortURL)
	api.Put("/update/:shortURL", controllers.UpdateShortURL)
	api.Get("/access/:shortURL", controllers.GetShortURLValidations)        // Step 1
	api.Get("/validate/:shortURL", controllers.ValidateAndRedirectShortURL) // Step 2

}
