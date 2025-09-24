/*
Copyright © 2025 NAME HERE <email>
*/

package tests

import (
	"github.com/dinoDanic/diny/ollama"
	"strings"
	"testing"
)

func TestMainStreamHelloCall(t *testing.T) {
	response, err := ollama.MainStream("hello")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if response == "" {
		t.Error("Expected non-empty response, but got empty string")
	}

	// Basic validation that we got some kind of response
	if len(strings.TrimSpace(response)) < 1 {
		t.Error("Expected meaningful response, but got whitespace only")
	}

	t.Logf("Ollama streaming response: %s", response)
}
