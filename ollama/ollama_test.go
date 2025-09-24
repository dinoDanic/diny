/*
Copyright © 2025 NAME HERE <email>
*/

package ollama

import (
	"strings"
	"testing"
)

func TestMainStreamHelloCall(t *testing.T) {
	response, err := MainStream("hello")

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

func TestMainHelloCall(t *testing.T) {
	// Test the non-streaming API call with "hello"
	response, err := Main("hello")

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

	t.Logf("Ollama response: %s", response)
}
