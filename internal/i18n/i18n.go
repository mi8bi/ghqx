package i18n

// Locale represents a language locale.
type Locale string

const (
	LocaleJA Locale = "ja"
	LocaleEN Locale = "en"
	// Add other locales as needed
)

// currentLocale holds the currently active locale.
var currentLocale Locale = LocaleJA // Default to Japanese

// messages holds all localized strings.
// Key: locale, Value: map of message keys to localized strings.
var messages = make(map[Locale]map[string]string)

func init() {
	// Initialize with Japanese messages by default
	loadJapaneseMessages()
}

// SetLocale sets the active locale.
func SetLocale(l Locale) {
	currentLocale = l
	// Load messages for the new locale if not already loaded
	if _, ok := messages[currentLocale]; !ok {
		// In a real app, you might load from files here
		// For now, we only have hardcoded ja and a placeholder en
		if currentLocale == LocaleEN {
			loadEnglishMessages()
		} else {
			loadJapaneseMessages() // Fallback
		}
	}
}

// T translates a message key into the current locale's string.
func T(key string) string {
	if bundle, ok := messages[currentLocale]; ok {
		if msg, ok := bundle[key]; ok {
			return msg
		}
	}
	// Fallback to English if not found in current locale or default
	if bundle, ok := messages[LocaleEN]; ok {
		if msg, ok := bundle[key]; ok {
			return msg
		}
	}
	// If all else fails, return the key itself (or a default error message)
	return "MISSING_TRANSLATION: " + key
}

// RegisterMessages adds a map of messages for a given locale.
func RegisterMessages(l Locale, msgs map[string]string) {
	if _, ok := messages[l]; !ok {
		messages[l] = make(map[string]string)
	}
	for k, v := range msgs {
		messages[l][k] = v
	}
}
