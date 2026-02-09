package sep24

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/stellar/sep-reference/reference/go/internal/db"
)

func (s *Service) handleGetTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	account := accountFromRequest(r)
	if account == "" {
		writeError(w, http.StatusForbidden, "missing subject")
		return
	}

	id := r.URL.Query().Get("id")
	externalID := r.URL.Query().Get("external_transaction_id")
	stellarID := r.URL.Query().Get("stellar_transaction_id")
	if id == "" && externalID == "" && stellarID == "" {
		writeError(w, http.StatusBadRequest, "missing transaction identifier")
		return
	}

	tx, ok := s.findTransactionForAccount(account, id, externalID, stellarID)
	if !ok {
		writeError(w, http.StatusNotFound, "transaction not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"transaction": s.toSEP24Transaction(tx)})
}

func (s *Service) handleListTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	account := accountFromRequest(r)
	if account == "" {
		writeError(w, http.StatusForbidden, "missing subject")
		return
	}

	limit := 10
	if raw := r.URL.Query().Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			writeError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = parsed
	}

	assetCode := strings.TrimSpace(r.URL.Query().Get("asset_code"))
	if assetCode != "" && !s.assetSupported(assetCode) {
		writeError(w, http.StatusBadRequest, "unsupported asset")
		return
	}

	kindFilter := strings.TrimSpace(r.URL.Query().Get("kind"))
	if kindFilter != "" && kindFilter != "deposit" && kindFilter != "withdrawal" {
		writeError(w, http.StatusBadRequest, "invalid kind")
		return
	}

	var noOlderThan time.Time
	if raw := strings.TrimSpace(r.URL.Query().Get("no_older_than")); raw != "" {
		parsed, err := time.Parse(time.RFC3339Nano, raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid no_older_than")
			return
		}
		noOlderThan = parsed
	}

	all := s.TxStore.ListByAccount(account, 0, "")
	sort.Slice(all, func(i, j int) bool {
		return all[i].StartedAt.After(all[j].StartedAt)
	})

	transactions := make([]map[string]any, 0, len(all))
	for _, tx := range all {
		mapped := s.toSEP24Transaction(tx)
		if assetCode != "" && tx.AssetCode != assetCode {
			continue
		}
		if kindFilter != "" && mapped["kind"] != kindFilter {
			continue
		}
		if !noOlderThan.IsZero() && tx.StartedAt.Before(noOlderThan) {
			continue
		}
		transactions = append(transactions, mapped)
	}

	if limit < len(transactions) {
		transactions = transactions[:limit]
	}
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
	writeJSON(w, http.StatusOK, map[string]float64{"fee": fee})
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

func (s *Service) findTransactionForAccount(account, id, externalID, stellarID string) (db.Transaction, bool) {
	for _, tx := range s.TxStore.ListByAccount(account, 0, "") {
		if id != "" && tx.ID == id {
			return tx, true
		}
		if externalID != "" && tx.ExternalTransactionID == externalID {
			return tx, true
		}
		if stellarID != "" && tx.StellarTransactionID == stellarID {
			return tx, true
		}
	}
	return db.Transaction{}, false
}

func (s *Service) toSEP24Transaction(tx db.Transaction) map[string]any {
	kind := "deposit"
	if tx.Kind == "withdraw" || tx.Kind == "withdrawal" {
		kind = "withdrawal"
	}

	status := tx.Status
	if status == "" {
		status = StatusIncomplete
	}

	moreInfoURL := tx.URL
	if strings.TrimSpace(moreInfoURL) == "" {
		moreInfoURL = fmt.Sprintf("http://%s/sep24/interactive/status?id=%s", s.Config.HomeDomain, tx.ID)
	}

	out := map[string]any{
		"id":            tx.ID,
		"kind":          kind,
		"status":        status,
		"more_info_url": moreInfoURL,
		"started_at":    tx.StartedAt,
		"updated_at":    tx.UpdatedAt,
		"asset_code":    tx.AssetCode,
	}
	if tx.Amount != "" {
		out["amount_in"] = tx.Amount
		out["amount_out"] = tx.Amount
		out["amount_in_asset"] = tx.AssetCode
		out["amount_out_asset"] = tx.AssetCode
	}
	if tx.StellarTransactionID != "" {
		out["stellar_transaction_id"] = tx.StellarTransactionID
	}
	if tx.ExternalTransactionID != "" {
		out["external_transaction_id"] = tx.ExternalTransactionID
	}

	if kind == "deposit" {
		out["to"] = tx.Account
		if tx.From != "" {
			out["from"] = tx.From
		}
		return out
	}

	out["from"] = tx.Account
	if tx.To != "" {
		out["to"] = tx.To
	}
	return out
}
