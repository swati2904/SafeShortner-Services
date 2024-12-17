package models

type Shortner struct {
	ID                   int      `json:"id" bson:"id" validate:"required"`
	LongURL              string   `json:"longURL" bson:"longURL" validate:"required"`
	ShortURL             string   `json:"shortURL" bson:"shortURL" validate:"required"`
	CheckedPasscode      bool     `json:"checkedPasscode,omitempty" bson:"checkedPasscode,omitempty"`
	Passcode             string   `json:"passcode,omitempty" bson:"passcode,omitempty"`
	CheckedCaptcha       bool     `json:"checkedCaptcha,omitempty" bson:"checkedCaptcha,omitempty"`
	Captcha              string   `json:"captcha,omitempty" bson:"captcha,omitempty"`
	CheckedAccessControl bool     `json:"checkedAccessControl,omitempty" bson:"checkedAccessControl,omitempty"`
	CountryBlacklist     []string `json:"countryBlacklist,omitempty" bson:"countryBlacklist,omitempty"`
	CheckedTimeZone      bool     `json:"checkedTimeZone,omitempty" bson:"checkedTimeZone,omitempty"`
	ExpiryDateTime       string   `json:"expiryDateTime,omitempty" bson:"expiryDateTime,omitempty"`
}
