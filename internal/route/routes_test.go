package route

import (
	"net/http"
	"testing"
)

func TestSessionOptions_LocalOriginsUseLaxCookie(t *testing.T) {
	options := sessionOptions("http://localhost:3000,http://127.0.0.1:5173")

	if options.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSiteLaxMode, got %v", options.SameSite)
	}

	if options.Secure {
		t.Fatal("expected Secure to be false for local origins")
	}
}

func TestSessionOptions_RemoteOriginsUseCrossSiteCookie(t *testing.T) {
	options := sessionOptions("https://whiteboard.example.com")

	if options.SameSite != http.SameSiteNoneMode {
		t.Fatalf("expected SameSiteNoneMode, got %v", options.SameSite)
	}

	if !options.Secure {
		t.Fatal("expected Secure to be true for remote origins")
	}
}

func TestSessionOptions_HTTPPrivateIPUsesLaxCookie(t *testing.T) {
	options := sessionOptions("http://172.20.10.10:3000")

	if options.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSiteLaxMode, got %v", options.SameSite)
	}

	if options.Secure {
		t.Fatal("expected Secure to be false for plain HTTP private IP origins")
	}
}
