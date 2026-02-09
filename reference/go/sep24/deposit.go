package sep24

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stellar/sep-reference/reference/go/internal/db"
)

func (s *Service) handleDepositInteractive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	defer r.Body.Close()

	var req InteractiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json request")
		return
	}

	account := accountFromRequest(r)
	if account == "" {
		writeError(w, http.StatusUnauthorized, "missing subject")
		return
	}
	if req.Account != "" && req.Account != account {
		writeError(w, http.StatusForbidden, "account mismatch")
		return
	}
	if req.AssetCode == "" {
		writeError(w, http.StatusBadRequest, "missing asset_code")
		return
	}
	if !s.assetSupported(req.AssetCode) {
		writeError(w, http.StatusBadRequest, "unsupported asset")
		return
	}

	now := s.Now()
	id := transactionID("dep", account, req.AssetCode, now)
	url := fmt.Sprintf("http://%s/sep24/interactive/deposit?id=%s", s.Config.HomeDomain, id)
	tx := db.Transaction{
		ID:        id,
		Kind:      "deposit",
		Status:    StatusIncomplete,
		Account:   account,
		AssetCode: req.AssetCode,
		Amount:    req.Amount,
		URL:       url,
		StartedAt: now,
		UpdatedAt: now,
		KYCFields: []string{"first_name", "last_name", "email_address"},
	}
	if err := s.TxStore.Create(tx); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create transaction")
		return
	}

	if s.CustomerStore != nil {
		_ = s.CustomerStore.Put(db.Customer{Account: account, Fields: map[string]string{}})
	}

	writeJSON(w, http.StatusOK, InteractiveResponse{
		ID:   id,
		Type: "interactive_customer_info_needed",
		URL:  url,
	})
}
