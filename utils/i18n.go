// Package utils provides common utility helpers including standardized JSON responses and i18n support.
package utils

import (
	"encoding/json"
	"fmt"

	"english-course-api/locales"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

// InitI18n initializes the i18n bundle and loads translation dictionary files from the embedded filesystem.
func InitI18n() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load dictionary files from embedded filesystem (locales.FS)
	if _, err := bundle.LoadMessageFileFS(locales.FS, "en.json"); err != nil {
		log.Warn().Err(err).Msg("Failed to load embedded locales/en.json")
	}
	if _, err := bundle.LoadMessageFileFS(locales.FS, "id.json"); err != nil {
		log.Warn().Err(err).Msg("Failed to load embedded locales/id.json")
	}

	log.Info().Msg("i18n Bundle successfully initialized with embedded locales (en, id) 🌐")
}

// Translate localizes a given messageID based on the active language stored in the Gin context.
func Translate(c *gin.Context, messageID string, templateData map[string]interface{}) string {
	if bundle == nil {
		return messageID
	}

	langStr := "en"
	if c != nil {
		if val, exists := c.Get("lang"); exists {
			if s, ok := val.(string); ok && s != "" {
				langStr = s
			}
		}
	}

	localizer := i18n.NewLocalizer(bundle, langStr)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
	if err != nil {
		return messageID
	}
	return msg
}

// TranslateValidationError translates validator/v10 field error tags into human-readable localized messages.
func TranslateValidationError(c *gin.Context, fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()

	msgKey := fmt.Sprintf("validation.%s", tag)
	templateData := map[string]interface{}{
		"Field": field,
		"Param": param,
		"Tag":   tag,
	}

	translated := Translate(c, msgKey, templateData)
	// If the specific key is not defined, fallback to generic validation message
	if translated == msgKey {
		return Translate(c, "validation.default", templateData)
	}

	return translated
}

