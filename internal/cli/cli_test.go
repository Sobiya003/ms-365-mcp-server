package cli

import "testing"

func TestParseCloudValidation(t *testing.T) {
	_, err := Parse([]string{"--cloud", "moon"})
	if err == nil {
		t.Fatal("expected invalid cloud error")
	}
}

func TestParsePresetValidation(t *testing.T) {
	_, err := Parse([]string{"--preset", "mail,unknown"})
	if err == nil {
		t.Fatal("expected invalid preset error")
	}
}
