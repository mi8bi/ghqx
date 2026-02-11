package i18n

import "testing"

func TestRegisterSetAndTranslate(t *testing.T) {
	// Backup state
	prevLocale := currentLocale
	prevMessages := make(map[Locale]map[string]string)
	for k, v := range messages {
		m := make(map[string]string)
		for kk, vv := range v {
			m[kk] = vv
		}
		prevMessages[k] = m
	}
	defer func() {
		currentLocale = prevLocale
		messages = prevMessages
	}()

	// Register English messages and switch
	RegisterMessages(LocaleEN, map[string]string{"key.hello": "hello"})
	SetLocale(LocaleEN)
	if got := T("key.hello"); got != "hello" {
		t.Fatalf("T returned %q, want %q", got, "hello")
	}

	// Missing key should fallback to english or return MISSING_TRANSLATION
	SetLocale(LocaleJA)
	RegisterMessages(LocaleEN, map[string]string{"only.en": "enval"})
	// ensure missing in ja but present in en
	if got := T("only.en"); got != "enval" {
		t.Fatalf("fallback to en failed, got %q", got)
	}
}
