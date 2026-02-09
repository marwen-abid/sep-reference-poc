package sep24

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (s *Service) handleGetTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id")
		return
	}

	tx, ok := s.TxStore.GetByID(id)
	if !ok {
		writeError(w, http.StatusNotFound, "transaction not found")
		return
	}
	if tx.Account != accountFromRequest(r) {
		writeError(w, http.StatusForbidden, "transaction does not belong to account")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"transaction": tx})
}

func (s *Service) handleListTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	limit := 10
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			limit = parsed
		}
	}
	cursor := r.URL.Query().Get("cursor")
	transactions := s.TxStore.ListByAccount(accountFromRequest(r), limit, cursor)
	writeJSON(w, http.StatusOK, map[string]any{"transactions": transactions})
}

func (s *Service) handleGetFee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	operation := r.URL.Query().Get("operation")
	assetCode := r.URL.Query().Get("asset_code")
	amountRaw := r.URL.Query().Get("amount")
	if operation == "" || assetCode == "" || amountRaw == "" {
		writeError(w, http.StatusBadRequest, "missing required query params")
		return
	}
	if operation != "deposit" && operation != "withdraw" {
		writeError(w, http.StatusBadRequest, "invalid operation")
		return
	}
	asset, ok := s.assetConfig(assetCode)
	if !ok {
		writeError(w, http.StatusBadRequest, "unsupported asset")
		return
	}
	amount, err := strconv.ParseFloat(amountRaw, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid amount")
		return
	}
	fee := asset.FeeFixed + (amount * asset.FeePercent / 100.0)
	writeJSON(w, http.StatusOK, map[string]string{"fee": formatAmount(fee)})
}

func (s *Service) assetSupported(code string) bool {
	_, ok := s.assetConfig(code)
	return ok
}

func (s *Service) assetConfig(code string) (asset struct {
	Code       string
	Enabled    bool
	FeeFixed   float64
	FeePercent float64
}, ok bool) {
	for _, a := range s.Config.Assets {
		if a.Code == code && a.Enabled {
			return struct {
				Code       string
				Enabled    bool
				FeeFixed   float64
				FeePercent float64
			}{a.Code, a.Enabled, a.FeeFixed, a.FeePercent}, true
		}
	}
	return asset, false
}

func formatAmount(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

func transactionID(prefix, account, asset string, t time.Time) string {
	h := sha1.Sum([]byte(fmt.Sprintf("%s|%s|%s|%d", prefix, account, asset, t.UnixNano())))
	return prefix + "-" + hex.EncodeToString(h[:8])
}
