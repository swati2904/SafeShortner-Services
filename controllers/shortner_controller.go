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
	"go.mongodb.org/mongo-driver/bson"
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

//update

func UpdateShortURL(c *fiber.Ctx) error {
	// Parse the request parameters
	shortURL := c.Params("shortURL")
	if shortURL == "" {
		return responses.JSONResponse(c, fiber.StatusBadRequest, "Short URL is required", nil)
	}

	// Parse the request body
	var input struct {
		NewShortURL          string   `json:"newShortURL,omitempty"`
		CheckedPasscode      bool     `json:"checkedPasscode,omitempty"`
		Passcode             string   `json:"passcode,omitempty"`
		CheckedCaptcha       bool     `json:"checkedCaptcha,omitempty"`
		Captcha              string   `json:"captcha,omitempty"`
		CheckedAccessControl bool     `json:"checkedAccessControl,omitempty"`
		CountryBlacklist     []string `json:"countryBlacklist,omitempty"`
		CheckedTimeZone      bool     `json:"checkedTimeZone,omitempty"`
		ExpiryDateTime       string   `json:"expiryDateTime,omitempty"`
	}
	if err := c.BodyParser(&input); err != nil {
		return responses.JSONResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validation for fields
	if input.NewShortURL != "" {
		// Check if the new short URL already exists
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		existingCount, err := configs.URLCollection.CountDocuments(ctx, bson.M{"shortURL": input.NewShortURL})
		if err != nil {
			return responses.JSONResponse(c, fiber.StatusInternalServerError, "Database query error", nil)
		}
		if existingCount > 0 {
			return responses.JSONResponse(c, fiber.StatusConflict, "Short URL already in use", nil)
		}
	}

	if input.CheckedPasscode && (len(input.Passcode) != 6 || input.Passcode == "") {
		return responses.JSONResponse(c, fiber.StatusBadRequest, "Passcode must be 6 digits long", nil)
	}

	if input.CheckedCaptcha && (input.Captcha != "cell" && input.Captcha != "text") {
		return responses.JSONResponse(c, fiber.StatusBadRequest, "Captcha must be either 'cell' or 'text'", nil)
	}

	if input.CheckedAccessControl && (len(input.CountryBlacklist) == 0) {
		return responses.JSONResponse(c, fiber.StatusBadRequest, "Country blacklist cannot be empty when Access Control is enabled", nil)
	}

	// Build the update document
	updateFields := bson.M{}
	if input.NewShortURL != "" {
		updateFields["shortURL"] = input.NewShortURL
	}
	if input.CheckedPasscode {
		updateFields["checkedPasscode"] = true
		updateFields["passcode"] = input.Passcode
	}
	if input.CheckedCaptcha {
		updateFields["checkedCaptcha"] = true
		updateFields["captcha"] = input.Captcha
	}
	if input.CheckedAccessControl {
		updateFields["checkedAccessControl"] = true
		updateFields["countryBlacklist"] = input.CountryBlacklist
	}
	if input.CheckedTimeZone {
		updateFields["checkedTimeZone"] = true
		updateFields["expiryDateTime"] = input.ExpiryDateTime
	}

	// Update the document in the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"shortURL": shortURL}
	update := bson.M{"$set": updateFields}

	result, err := configs.URLCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return responses.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update the short URL", nil)
	}
	if result.MatchedCount == 0 {
		return responses.JSONResponse(c, fiber.StatusNotFound, "Short URL not found", nil)
	}

	return responses.JSONResponse(c, fiber.StatusOK, "Short URL updated successfully", updateFields)
}
