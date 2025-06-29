/*
Copyright (c) 2025 Regi Ellis

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/


package civitai

import (
	"strings"
	"testing"
)

func TestSecureTokenMethods(t *testing.T) {
	t.Run("Client without token", func(t *testing.T) {
		client := NewClientWithoutAuth()
		
		// Test HasAPIToken
		if client.HasAPIToken() {
			t.Error("Expected HasAPIToken to return false for client without token")
		}
		
		// Test IsAuthenticated
		if client.IsAuthenticated() {
			t.Error("Expected IsAuthenticated to return false for client without token")
		}
		
		// Test GetMaskedAPIToken
		masked := client.GetMaskedAPIToken()
		if masked != "none" {
			t.Errorf("Expected 'none' for empty token, got '%s'", masked)
		}
	})

	t.Run("Client with token", func(t *testing.T) {
		token := "sk-1234567890abcdef1234567890abcdef"
		client := NewClient(token)
		
		// Test HasAPIToken
		if !client.HasAPIToken() {
			t.Error("Expected HasAPIToken to return true for client with token")
		}
		
		// Test IsAuthenticated
		if !client.IsAuthenticated() {
			t.Error("Expected IsAuthenticated to return true for client with token")
		}
		
		// Test GetMaskedAPIToken
		masked := client.GetMaskedAPIToken()
		expected := "sk-12345" + strings.Repeat("*", len(token)-8)
		if masked != expected {
			t.Errorf("Expected '%s', got '%s'", expected, masked)
		}
		
		// Verify original token is not exposed in masked version
		if strings.Contains(masked, "abcdef") {
			t.Error("Masked token should not contain original token parts")
		}
		
		// Verify GetAPIToken still works (backward compatibility)
		if client.GetAPIToken() != token {
			t.Errorf("Expected original token '%s', got '%s'", token, client.GetAPIToken())
		}
	})

	t.Run("Short token masking", func(t *testing.T) {
		shortToken := "sk-123"
		client := NewClient(shortToken)
		
		masked := client.GetMaskedAPIToken()
		expected := strings.Repeat("*", len(shortToken))
		
		if masked != expected {
			t.Errorf("Expected '%s' for short token, got '%s'", expected, masked)
		}
	})

	t.Run("Very short token masking", func(t *testing.T) {
		veryShortToken := "abc"
		client := NewClient(veryShortToken)
		
		masked := client.GetMaskedAPIToken()
		expected := "***"
		
		if masked != expected {
			t.Errorf("Expected '%s' for very short token, got '%s'", expected, masked)
		}
	})

	t.Run("Token masking preserves prefix", func(t *testing.T) {
		token := "api-key-1234567890abcdef"
		client := NewClient(token)
		
		masked := client.GetMaskedAPIToken()
		
		// Should start with first 8 characters
		if !strings.HasPrefix(masked, "api-key-") {
			t.Errorf("Expected masked token to start with 'api-key-', got '%s'", masked)
		}
		
		// Should end with asterisks
		asteriskCount := strings.Count(masked, "*")
		expectedAsterisks := len(token) - 8
		if asteriskCount != expectedAsterisks {
			t.Errorf("Expected %d asterisks, got %d", expectedAsterisks, asteriskCount)
		}
	})

	t.Run("Security warning in GetAPIToken", func(t *testing.T) {
		// This test documents the security concern with GetAPIToken
		// In real usage, developers should prefer HasAPIToken() or GetMaskedAPIToken()
		token := "sensitive-token-12345"
		client := NewClient(token)
		
		// GetAPIToken exposes the full token (security risk)
		exposedToken := client.GetAPIToken()
		if exposedToken != token {
			t.Errorf("GetAPIToken should return full token for backward compatibility")
		}
		
		// Safer alternatives
		if !client.HasAPIToken() {
			t.Error("HasAPIToken should return true")
		}
		
		masked := client.GetMaskedAPIToken()
		if strings.Contains(masked, "12345") {
			t.Error("Masked token should not contain sensitive parts")
		}
	})
}

func TestTokenSecurityBestPractices(t *testing.T) {
	t.Run("Recommend using HasAPIToken instead of GetAPIToken", func(t *testing.T) {
		client := NewClient("my-secret-token")
		
		// GOOD: Check if token exists without exposing it
		if client.HasAPIToken() {
			// Token is configured, proceed with authenticated requests
		}
		
		// BAD: Exposes token in logs/memory
		// token := client.GetAPIToken()
		// log.Printf("Using token: %s", token) // Security risk!
		
		// GOOD: Use masked version for logging
		masked := client.GetMaskedAPIToken()
		if !strings.Contains(masked, "*") {
			t.Error("Masked token should contain asterisks for security")
		}
	})

	t.Run("IsAuthenticated is alias for HasAPIToken", func(t *testing.T) {
		clientWithToken := NewClient("token")
		clientWithoutToken := NewClientWithoutAuth()
		
		// Both methods should return the same result
		if clientWithToken.HasAPIToken() != clientWithToken.IsAuthenticated() {
			t.Error("HasAPIToken and IsAuthenticated should return same value")
		}
		
		if clientWithoutToken.HasAPIToken() != clientWithoutToken.IsAuthenticated() {
			t.Error("HasAPIToken and IsAuthenticated should return same value")
		}
	})
}

func TestMaskedTokenLogging(t *testing.T) {
	t.Run("Safe logging example", func(t *testing.T) {
		client := NewClient("sk-1234567890abcdefghijklmnop")
		
		// Safe for logging
		masked := client.GetMaskedAPIToken()
		
		// Verify it's safe to log
		if strings.Contains(masked, "ghijklmnop") {
			t.Error("Masked token contains sensitive suffix")
		}
		
		if !strings.HasPrefix(masked, "sk-12345") {
			t.Error("Masked token should preserve recognizable prefix")
		}
		
		// Count asterisks
		asterisks := strings.Count(masked, "*")
		if asterisks == 0 {
			t.Error("Masked token should contain asterisks")
		}
	})
}
