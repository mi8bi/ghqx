package main

import (
	"os"
	"testing"

	"github.com/mi8bi/ghqx/internal/i18n"
)

func TestParseLocaleFromEnvCases(t *testing.T) {
	if parseLocaleFromEnv("en") != i18n.LocaleEN {
		t.Fatalf("en should map to LocaleEN")
	}
	if parseLocaleFromEnv("ja_JP") != i18n.LocaleJA {
		t.Fatalf("ja_JP should map to LocaleJA")
	}
	if parseLocaleFromEnv("xx") != "" {
		t.Fatalf("unknown should return empty")
	}
}

func TestGetOSLanguageLocalePrecedence(t *testing.T) {
	origLC := os.Getenv("LC_ALL")
	origLANG := os.Getenv("LANG")
	origLANGUAGE := os.Getenv("LANGUAGE")
	defer os.Setenv("LC_ALL", origLC)
	defer os.Setenv("LANG", origLANG)
	defer os.Setenv("LANGUAGE", origLANGUAGE)

	os.Unsetenv("LC_ALL")
	os.Unsetenv("LANG")
	os.Unsetenv("LANGUAGE")

	// default should be JA
	if getOSLanguageLocale() != i18n.LocaleJA {
		t.Fatalf("default locale should be JA")
	}

	os.Setenv("LANG", "en_US.UTF-8")
	if getOSLanguageLocale() != i18n.LocaleEN {
		t.Fatalf("LANG en_US should map to EN")
	}

	os.Unsetenv("LANG")
	os.Setenv("LC_ALL", "ja_JP.UTF-8")
	if getOSLanguageLocale() != i18n.LocaleJA {
		t.Fatalf("LC_ALL ja_JP should map to JA")
	}

	os.Unsetenv("LC_ALL")
	os.Setenv("LANGUAGE", "en")
	if getOSLanguageLocale() != i18n.LocaleEN {
		t.Fatalf("LANGUAGE en should map to EN")
	}
}
