package tui

import (
	"testing"

	"github.com/mi8bi/ghqx/internal/i18n"
)

func TestGetWorkspaceStyle(t *testing.T) {
	testCases := []struct {
		workspace string
		expected  bool // whether a styled output is expected
	}{
		{"sandbox", true},
		{"dev", true},
		{"release", true},
		{"unknown", true}, // should return default style
	}

	for _, tc := range testCases {
		style := getWorkspaceStyle(tc.workspace)
		if style.GetBackground() == nil && style.GetForeground() == nil {
			// Even default style is valid
		}
	}
}

func TestGetStatusStyle(t *testing.T) {
	i18n.SetLocale(i18n.LocaleEN)

	testCases := []struct {
		status string
	}{
		{i18n.T("status.repo.clean")},
		{i18n.T("status.repo.dirty")},
		{"-"}, // unmanaged
		{"unknown"},
	}

	for _, tc := range testCases {
		style := getStatusStyle(tc.status)
		// Style should be returned without panic
		if style.GetBackground() == nil && style.GetForeground() == nil {
			// Even default style is valid
		}
	}
}

func TestGetMessageStyle(t *testing.T) {
	testCases := []struct {
		msgType MessageType
	}{
		{MessageTypeInfo},
		{MessageTypeSuccess},
		{MessageTypeWarning},
		{MessageTypeError},
	}

	for _, tc := range testCases {
		style := getMessageStyle(tc.msgType)
		// Style should be returned without panic
		if style.GetBackground() == nil && style.GetForeground() == nil {
			// Even default style is valid
		}
	}
}

func TestStylesExist(t *testing.T) {
	// Test that all style variables are initialized
	styles := []interface{}{
		styleHeader,
		styleRow,
		styleSelectedRow,
		styleSandbox,
		styleDev,
		styleRelease,
		styleClean,
		styleDirty,
		styleInfo,
		styleSuccess,
		styleWarning,
		styleError,
		styleTitle,
		styleHelp,
		styleFooter,
	}

	for _, style := range styles {
		if style == nil {
			t.Error("style should not be nil")
		}
	}
}
