package sep10

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBuildVerifyIssueToken(t *testing.T) {
	service := NewService(
		"GSERVERACCOUNTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"server-secret",
		"jwt-secret",
		"localhost:8080",
		"localhost:8080",
		"Test SDF Network ; September 2015",
		5*time.Minute,
		15*time.Minute,
	)

	challenge, err := service.BuildChallenge("GCLIENTACCOUNTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "wallet.example", "", "")
	if err != nil {
		t.Fatalf("build challenge: %v", err)
	}

	signed, err := AddClientSignature(challenge, "client-secret")
	if err != nil {
		t.Fatalf("add client signature: %v", err)
	}

	token, err := service.VerifyAndIssueToken(signed)
	if err != nil {
		t.Fatalf("verify and issue token: %v", err)
	}

	claims, err := VerifyToken(token, "jwt-secret")
	if err != nil {
		t.Fatalf("verify token: %v", err)
	}
	if claims.Subject == "" {
		t.Fatalf("expected subject in claims")
	}
}

func TestRejectMissingClientSignature(t *testing.T) {
	service := NewService(
		"GSERVERACCOUNTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"server-secret",
		"jwt-secret",
		"localhost:8080",
		"localhost:8080",
		"Test SDF Network ; September 2015",
		5*time.Minute,
		15*time.Minute,
	)

	challenge, err := service.BuildChallenge("GCLIENTACCOUNTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "", "", "")
	if err != nil {
		t.Fatalf("build challenge: %v", err)
	}

	if _, err := service.VerifyAndIssueToken(challenge); err == nil {
		t.Fatalf("expected missing client signature error")
	}
}

func TestHTTPAuthEndpoints(t *testing.T) {
	service := NewService(
		"GSERVERACCOUNTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"server-secret",
		"jwt-secret",
		"localhost:8080",
		"localhost:8080",
		"Test SDF Network ; September 2015",
		5*time.Minute,
		15*time.Minute,
	)
	handler := NewHTTPHandler(service)

	challengeReq := httptest.NewRequest(http.MethodGet, "/auth?account=GCLIENTACCOUNTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", nil)
	challengeRec := httptest.NewRecorder()
	handler.ServeHTTP(challengeRec, challengeReq)
	if challengeRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", challengeRec.Code)
	}

	var challengeBody map[string]string
	if err := json.Unmarshal(challengeRec.Body.Bytes(), &challengeBody); err != nil {
		t.Fatalf("decode challenge response: %v", err)
	}
	signed, err := AddClientSignature(challengeBody["transaction"], "client-secret")
	if err != nil {
		t.Fatalf("add client signature: %v", err)
	}

	payload, _ := json.Marshal(map[string]string{"transaction": signed})
	tokenReq := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(payload))
	tokenRec := httptest.NewRecorder()
	handler.ServeHTTP(tokenRec, tokenReq)
	if tokenRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", tokenRec.Code, tokenRec.Body.String())
	}
}
