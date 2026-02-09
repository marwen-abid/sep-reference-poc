package sep10

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

func TestBuildVerifyIssueToken(t *testing.T) {
	serverKP := mustRandomKeypair(t)
	clientKP := mustRandomKeypair(t)

	service := NewService(
		serverKP.Address(),
		serverKP.Seed(),
		"jwt-secret",
		"localhost:8080",
		"localhost:8080",
		DefaultNetworkPassphrase,
		5*time.Minute,
		15*time.Minute,
	)
	service.AccountSigners = staticSignerLoader(clientKP.Address())

	challenge, err := service.BuildChallenge(clientKP.Address(), "", "", "")
	if err != nil {
		t.Fatalf("build challenge: %v", err)
	}

	signed, err := AddClientSignatureWithNetworkPassphrase(challenge, clientKP.Seed(), DefaultNetworkPassphrase)
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
	if claims.Subject != clientKP.Address() {
		t.Fatalf("unexpected subject: %s", claims.Subject)
	}
}

func TestRejectMissingClientSignature(t *testing.T) {
	serverKP := mustRandomKeypair(t)
	clientKP := mustRandomKeypair(t)

	service := NewService(
		serverKP.Address(),
		serverKP.Seed(),
		"jwt-secret",
		"localhost:8080",
		"localhost:8080",
		DefaultNetworkPassphrase,
		5*time.Minute,
		15*time.Minute,
	)
	service.AccountSigners = staticSignerLoader(clientKP.Address())

	challenge, err := service.BuildChallenge(clientKP.Address(), "", "", "")
	if err != nil {
		t.Fatalf("build challenge: %v", err)
	}

	if _, err := service.VerifyAndIssueToken(challenge); err == nil {
		t.Fatalf("expected missing client signature error")
	}
}

func TestHTTPAuthEndpoints(t *testing.T) {
	serverKP := mustRandomKeypair(t)
	clientKP := mustRandomKeypair(t)

	service := NewService(
		serverKP.Address(),
		serverKP.Seed(),
		"jwt-secret",
		"localhost:8080",
		"localhost:8080",
		DefaultNetworkPassphrase,
		5*time.Minute,
		15*time.Minute,
	)
	service.AccountSigners = staticSignerLoader(clientKP.Address())

	handler := NewHTTPHandler(service)

	challengeReq := httptest.NewRequest(http.MethodGet, "/auth?account="+clientKP.Address(), nil)
	challengeRec := httptest.NewRecorder()
	handler.ServeHTTP(challengeRec, challengeReq)
	if challengeRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", challengeRec.Code, challengeRec.Body.String())
	}

	var challengeBody map[string]string
	if err := json.Unmarshal(challengeRec.Body.Bytes(), &challengeBody); err != nil {
		t.Fatalf("decode challenge response: %v", err)
	}

	signed, err := AddClientSignatureWithNetworkPassphrase(
		challengeBody["transaction"],
		clientKP.Seed(),
		challengeBody["network_passphrase"],
	)
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

func mustRandomKeypair(t *testing.T) *keypair.Full {
	t.Helper()
	kp, err := keypair.Random()
	if err != nil {
		t.Fatalf("generate keypair: %v", err)
	}
	return kp
}

func staticSignerLoader(accountID string) AccountSignerLoader {
	return func(requested string) (txnbuild.SignerSummary, txnbuild.Threshold, bool, error) {
		if requested != accountID {
			return nil, 0, false, nil
		}
		return txnbuild.SignerSummary{
			accountID: 1,
		}, 1, true, nil
	}
}
