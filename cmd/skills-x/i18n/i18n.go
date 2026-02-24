package i18n

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

//go:embed locales/*.yaml
var localesFS embed.FS

var (
	currentLang  = "zh"
	messages     map[string]string
	messagesLock sync.RWMutex
	initialized  bool
)

// Init initializes the i18n package with the language from environment variable.
// Supported languages: zh (Chinese), en (English)
// Environment variable: LANG or SKILLS_LANG
func Init() error {
	lang := detectLanguage()
	return SetLanguage(lang)
}

// detectLanguage detects language from environment variables
func detectLanguage() string {
	// Priority: SKILLS_LANG > LANG > LC_ALL > default(zh)
	if lang := os.Getenv("SKILLS_LANG"); lang != "" {
		return normalizeLanguage(lang)
	}
	if lang := os.Getenv("LANG"); lang != "" {
		return normalizeLanguage(lang)
	}
	if lang := os.Getenv("LC_ALL"); lang != "" {
		return normalizeLanguage(lang)
	}
	return "zh"
}

// normalizeLanguage converts locale string to supported language code
func normalizeLanguage(locale string) string {
	locale = strings.ToLower(locale)
	if strings.HasPrefix(locale, "zh") {
		return "zh"
	}
	if strings.HasPrefix(locale, "en") {
		return "en"
	}
	// Default to Chinese for unsupported languages
	return "zh"
}

// SetLanguage sets the current language and loads corresponding messages
func SetLanguage(lang string) error {
	messagesLock.Lock()
	defer messagesLock.Unlock()

	if lang != "zh" && lang != "en" {
		lang = "zh"
	}

	data, err := localesFS.ReadFile(fmt.Sprintf("locales/%s.yaml", lang))
	if err != nil {
		return fmt.Errorf("failed to load language file %s: %w", lang, err)
	}

	var msgs map[string]string
	if err := yaml.Unmarshal(data, &msgs); err != nil {
		return fmt.Errorf("failed to parse language file %s: %w", lang, err)
	}

	currentLang = lang
	messages = msgs
	initialized = true
	return nil
}

// GetLanguage returns the current language code
func GetLanguage() string {
	messagesLock.RLock()
	defer messagesLock.RUnlock()
	return currentLang
}

// T translates a message key to the current language.
// If the key is not found, it returns the key itself.
func T(key string) string {
	messagesLock.RLock()
	defer messagesLock.RUnlock()

	if !initialized {
		return key
	}

	if msg, ok := messages[key]; ok {
		return msg
	}
	return key
}

// Tf translates a message key with format arguments
func Tf(key string, args ...interface{}) string {
	return fmt.Sprintf(T(key), args...)
}

// MustInit initializes the i18n package and panics on error
func MustInit() {
	if err := Init(); err != nil {
		panic(err)
	}
}
