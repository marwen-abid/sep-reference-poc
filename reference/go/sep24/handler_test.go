package sep24

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stellar/sep-reference/reference/go/internal/config"
	"github.com/stellar/sep-reference/reference/go/internal/db"
	"github.com/stellar/sep-reference/reference/go/internal/middleware"
	"github.com/stellar/sep-reference/reference/go/sep10"
)

func TestTransitionValidation(t *testing.T) {
	if err := ValidateTransition(StatusIncomplete, StatusPendingUserTransferStart); err != nil {
		t.Fatalf("expected valid transition, got %v", err)
	}
	if err := ValidateTransition(StatusIncomplete, StatusCompleted); err == nil {
		t.Fatalf("expected invalid transition error")
	}
}

func TestAuthRequiredOnProtectedRoutes(t *testing.T) {
	service, mux := testServiceAndMux()
	_ = service

	req := httptest.NewRequest(http.MethodPost, "/sep24/transactions/deposit/interactive", bytes.NewBufferString(`{"asset_code":"USDC"}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestDepositAndGetTransaction(t *testing.T) {
	service, mux := testServiceAndMux()
	token, err := sep10.IssueToken("GTESTACCOUNTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "localhost:8080", "", "localhost:8080", "jwt-secret", time.Now().UTC(), 10*time.Minute)
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}

	payload, _ := json.Marshal(InteractiveRequest{AssetCode: "USDC", Amount: "100.00"})
	req := httptest.NewRequest(http.MethodPost, "/sep24/transactions/deposit/interactive", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}

	var interactive InteractiveResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &interactive); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if _, ok := service.TxStore.GetByID(interactive.ID); !ok {
		t.Fatalf("expected transaction to be persisted")
	}

	getReq := httptest.NewRequest(http.MethodGet, "/sep24/transaction?id="+interactive.ID, nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", getRec.Code, getRec.Body.String())
	}
}

func testServiceAndMux() (*Service, *http.ServeMux) {
	cfg := config.Config{
		HomeDomain: "localhost:8080",
		JWTSecret:  "jwt-secret",
		Assets: []config.Asset{
			{Code: "USDC", Enabled: true, FeeFixed: 1.0, FeePercent: 0.1},
		},
	}
	txStore := db.NewMemoryTransactionStore()
	customerStore := db.NewMemoryCustomerStore()
	service := NewService(cfg, txStore, customerStore)

	mux := http.NewServeMux()
	service.RegisterRoutes(mux, middleware.SEP10Auth(cfg.JWTSecret))
	return service, mux
}
