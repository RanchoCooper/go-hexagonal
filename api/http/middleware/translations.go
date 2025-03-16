package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// Translations returns a middleware that handles translations for validation errors
func Translations() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize translators
		uni := ut.New(en.New(), zh.New())

		// Get locale from header or default to en
		locale := c.GetHeader("Accept-Language")
		if locale == "" {
			locale = "en"
		}

		// Get translator for the locale
		trans, _ := uni.GetTranslator(locale)

		// Register translator with validator
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			switch locale {
			case "zh":
				zh_translations.RegisterDefaultTranslations(v, trans)
			default:
				en_translations.RegisterDefaultTranslations(v, trans)
			}

			// Store translator in context
			c.Set("trans", trans)
		}

		c.Next()
	}
}
