package i18n

import (
	"embed"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var localeFS embed.FS

var bundle *i18n.Bundle

func Init() error {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, file := range []string{
		"locales/active.en.toml",
		"locales/active.zh.toml",
	} {
		buf, err := localeFS.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading locale %s: %w", file, err)
		}
		if _, err := bundle.ParseMessageFileBytes(buf, file); err != nil {
			return fmt.Errorf("parsing locale %s: %w", file, err)
		}
	}
	return nil
}

// T translates a message key using the locale from gin.Context.
func T(c *gin.Context, key string) string {
	return localize(c, key, nil)
}

// TWithData translates a message key with template data.
func TWithData(c *gin.Context, key string, data map[string]interface{}) string {
	return localize(c, key, data)
}

func localize(c *gin.Context, key string, data map[string]interface{}) string {
	lang := "en"
	if l, exists := c.Get("locale"); exists {
		if s, ok := l.(string); ok {
			lang = s
		}
	}
	localizer := i18n.NewLocalizer(bundle, lang)
	cfg := &i18n.LocalizeConfig{MessageID: key}
	if data != nil {
		cfg.TemplateData = data
	}
	msg, err := localizer.Localize(cfg)
	if err != nil {
		return key
	}
	return msg
}

// Error is a translatable error that carries an i18n key.
type Error struct {
	Key  string
	Data map[string]interface{}
}

func (e *Error) Error() string {
	return e.Key
}

// NewError creates a translatable error with the given i18n key.
func NewError(key string) *Error {
	return &Error{Key: key}
}

// TranslateError translates an error for user-facing display.
func TranslateError(c *gin.Context, err error) string {
	var te *Error
	if errors.As(err, &te) {
		if te.Data != nil {
			return TWithData(c, te.Key, te.Data)
		}
		return T(c, te.Key)
	}

	var nf *models.NotFoundError
	if errors.As(err, &nf) {
		translated := T(c, nf.Msg)
		if translated != nf.Msg {
			return translated
		}
		return nf.Msg
	}

	return err.Error()
}
