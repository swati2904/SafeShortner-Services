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

// Get fetch required validations
func GetShortURLValidations(c *fiber.Ctx) error {
	// Fetch the short URL from the request params
	shortURL := c.Params("shortURL")
	if shortURL == "" {
		return responses.JSONResponse(c, fiber.StatusBadRequest, "Short URL is required", nil)
	}

	// Fetch the document from the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var shortner models.Shortner
	err := configs.URLCollection.FindOne(ctx, bson.M{"shortURL": shortURL}).Decode(&shortner)
	if err != nil {
		return responses.JSONResponse(c, fiber.StatusNotFound, "Short URL not found", nil)
	}

	// Build the response to inform required validations
	requiredFields := bson.M{}
	if shortner.CheckedPasscode {
		requiredFields["passcodeRequired"] = true
	}
	if shortner.CheckedCaptcha {
		requiredFields["captchaRequired"] = true
	}
	if shortner.CheckedAccessControl {
		requiredFields["accessControlRequired"] = true
	}
	if shortner.CheckedTimeZone {
		requiredFields["timeZoneRequired"] = true
	}

	// Return required validation details
	return responses.JSONResponse(c, fiber.StatusOK, "Validations required for this short URL", requiredFields)
}

//get validate and redirect

// Step 2: GET /validate/:shortURL
func ValidateAndRedirectShortURL(c *fiber.Ctx) error {
	shortURL := c.Params("shortURL")
	if shortURL == "" {
		return responses.JSONResponse(c, fiber.StatusBadRequest, "Short URL is required", nil)
	}

	passcode := c.Query("passcode")
	captcha := c.Query("captcha")
	country := c.Query("country")
	currentTime := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var shortner models.Shortner
	err := configs.URLCollection.FindOne(ctx, bson.M{"shortURL": shortURL}).Decode(&shortner)
	if err != nil {
		return responses.JSONResponse(c, fiber.StatusNotFound, "Short URL not found", nil)
	}

	// Validate passcode
	if shortner.CheckedPasscode && shortner.Passcode != passcode {
		return responses.JSONResponse(c, fiber.StatusForbidden, "Invalid passcode", nil)
	}

	// Validate captcha
	if shortner.CheckedCaptcha && shortner.Captcha != captcha {
		return responses.JSONResponse(c, fiber.StatusForbidden, "Invalid captcha", nil)
	}

	// Access control
	if shortner.CheckedAccessControl && contains(shortner.CountryBlacklist, country) {
		return responses.JSONResponse(c, fiber.StatusForbidden, "Access from this country is restricted", nil)
	}

	// Expiry check
	if shortner.CheckedTimeZone {
		expiryTime, err := time.Parse("01/02/2006 03:04PM", shortner.ExpiryDateTime)
		if err != nil {
			return responses.JSONResponse(c, fiber.StatusInternalServerError, "Invalid expiry time format", nil)
		}
		if currentTime.After(expiryTime) {
			return responses.JSONResponse(c, fiber.StatusGone, "Link has expired", nil)
		}
	}

	// If all validations pass, redirect to the long URL
	return c.Redirect(shortner.LongURL)
}

// helper function to check if string exist
func contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
