package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Translator handles internationalization
type Translator struct {
	language     string
	translations map[string]interface{}
}

var globalTranslator *Translator

// Init initializes the global translator with the specified language
func Init(language string) error {
	t := &Translator{
		language: language,
	}

	if err := t.loadTranslations(); err != nil {
		// Fallback to English if loading fails
		t.language = "en"
		if err := t.loadTranslations(); err != nil {
			return fmt.Errorf("failed to load translations: %w", err)
		}
	}

	globalTranslator = t
	return nil
}

// loadTranslations loads translation file for current language
func (t *Translator) loadTranslations() error {
	// Get executable directory
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)

	// Try locales directory relative to executable
	localesPath := filepath.Join(exeDir, "locales", t.language+".json")

	// If not found, try relative to working directory (for development)
	if _, err := os.Stat(localesPath); os.IsNotExist(err) {
		localesPath = filepath.Join("locales", t.language+".json")
	}

	data, err := os.ReadFile(localesPath)
	if err != nil {
		return fmt.Errorf("failed to read translation file: %w", err)
	}

	var translations map[string]interface{}
	if err := json.Unmarshal(data, &translations); err != nil {
		return fmt.Errorf("failed to parse translation file: %w", err)
	}

	t.translations = translations
	return nil
}

// T translates a key using dot notation (e.g., "home.title")
func T(key string) string {
	if globalTranslator == nil {
		return key
	}
	return globalTranslator.get(key)
}

// get retrieves a translation by key using dot notation
func (t *Translator) get(key string) string {
	parts := strings.Split(key, ".")
	current := t.translations

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - should be a string
			if val, ok := current[part].(string); ok {
				return val
			}
			return key // Return key if not found
		}

		// Navigate deeper
		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			return key // Return key if path not found
		}
	}

	return key
}

// GetLanguage returns current language
func GetLanguage() string {
	if globalTranslator == nil {
		return "en"
	}
	return globalTranslator.language
}

// SetLanguage changes the current language and reloads translations
func SetLanguage(language string) error {
	return Init(language)
}
