package controllers

import (
	"SafeShotner-Services/configs"
	"SafeShotner-Services/models"
	"SafeShotner-Services/responses"
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Generate a random short URL
func generateShortURL() string {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Error generating short URL: %v", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

// Create a short URL
func CreateShortURL(c *fiber.Ctx) error {
	// Parse request body
	var input struct {
		LongURL string `json:"longURL" validate:"required"`
	}
	if err := c.BodyParser(&input); err != nil {
		return responses.JSONResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validate input
	if input.LongURL == "" {
		return responses.JSONResponse(c, fiber.StatusBadRequest, "Long URL is required", nil)
	}

	// Create a new URL object
	newURL := models.Shortner{
		ID:       int(time.Now().UnixNano() / 1e6), // Unique ID
		LongURL:  input.LongURL,
		ShortURL: generateShortURL(),
	}

	// Save to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := configs.URLCollection.InsertOne(ctx, newURL) // Use shared collection
	if err != nil {
		log.Printf("Database insert error: %v", err)
		return responses.JSONResponse(c, fiber.StatusInternalServerError, "Failed to save URL", nil)
	}

	return responses.JSONResponse(c, fiber.StatusCreated, "Short URL created successfully", newURL)
}
