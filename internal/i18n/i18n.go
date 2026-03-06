package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

//go:embed locales/*.json
var localesFS embed.FS

var (
	mu          sync.RWMutex
	currentLang string
	messages    map[string]string
	fallback    map[string]string
)

// Init loads the locale. If lang is empty, it detects from $LANG/$LC_ALL.
// Fallback is always "en".
func Init(lang string) {
	if lang == "" {
		lang = detectLang()
	}

	mu.Lock()
	defer mu.Unlock()

	fallback = loadLocale("en")
	if lang == "en" {
		messages = fallback
	} else {
		messages = loadLocale(lang)
	}
	currentLang = lang
}

// T returns the translated string for the given key.
func T(key string) string {
	mu.RLock()
	defer mu.RUnlock()

	if messages == nil {
		return key
	}
	if v, ok := messages[key]; ok {
		return v
	}
	if v, ok := fallback[key]; ok {
		return v
	}
	return key
}

// TF returns the translated string formatted with args (fmt.Sprintf).
func TF(key string, args ...any) string {
	return fmt.Sprintf(T(key), args...)
}

// Lang returns the current language code.
func Lang() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentLang
}

// SetLang switches the active language at runtime.
func SetLang(lang string) {
	mu.Lock()
	defer mu.Unlock()

	if lang == currentLang {
		return
	}
	if lang == "en" {
		messages = fallback
	} else {
		messages = loadLocale(lang)
	}
	currentLang = lang
}

func loadLocale(lang string) map[string]string {
	data, err := localesFS.ReadFile("locales/" + lang + ".json")
	if err != nil {
		return make(map[string]string)
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return make(map[string]string)
	}
	return m
}

func detectLang() string {
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		v := os.Getenv(env)
		if v == "" || v == "C" || v == "POSIX" {
			continue
		}
		// e.g. "pt_BR.UTF-8" → "pt-BR"
		v = strings.SplitN(v, ".", 2)[0]
		v = strings.ReplaceAll(v, "_", "-")
		return v
	}
	return "en"
}
