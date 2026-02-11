package configtui

import "testing"

func TestBoolStringIntHelpers(t *testing.T) {
	if boolToString(true) != "true" {
		t.Fatalf("boolToString(true) != \"true\"")
	}
	if boolToString(false) != "false" {
		t.Fatalf("boolToString(false) != \"false\"")
	}

	for _, s := range []string{"true", "yes", "1"} {
		if !stringToBool(s) {
			t.Fatalf("stringToBool(%q) should be true", s)
		}
	}
	if stringToBool("no") {
		t.Fatalf("stringToBool(\"no\") should be false")
	}

	if intToString(123) != "123" {
		t.Fatalf("intToString(123) != \"123\"")
	}

	if stringToInt("45") != 45 {
		t.Fatalf("stringToInt(\"45\") != 45")
	}
	if stringToInt("bad") != 0 {
		t.Fatalf("stringToInt(\"bad\") should return 0 for invalid input")
	}
}
