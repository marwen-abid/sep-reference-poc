package sep24

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/stellar/sep-reference/reference/go/internal/config"
	"github.com/stellar/sep-reference/reference/go/internal/db"
	"github.com/stellar/sep-reference/reference/go/internal/middleware"
)

type Service struct {
	Config        config.Config
	TxStore       db.TransactionStore
	CustomerStore db.CustomerStore
	Now           func() time.Time
}

type InteractiveRequest struct {
	AssetCode string `json:"asset_code"`
	Account   string `json:"account"`
	Amount    string `json:"amount"`
	Lang      string `json:"lang,omitempty"`
}

type InteractiveResponse struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

func NewService(cfg config.Config, txStore db.TransactionStore, customerStore db.CustomerStore) *Service {
	return &Service{
		Config:        cfg,
		TxStore:       txStore,
		CustomerStore: customerStore,
		Now:           func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.HandleFunc("/sep24/info", s.handleInfo)
	mux.Handle("/sep24/transactions/deposit/interactive", authMiddleware(http.HandlerFunc(s.handleDepositInteractive)))
	mux.Handle("/sep24/transactions/withdraw/interactive", authMiddleware(http.HandlerFunc(s.handleWithdrawInteractive)))
	mux.Handle("/sep24/transaction", authMiddleware(http.HandlerFunc(s.handleGetTransaction)))
	mux.Handle("/sep24/transactions", authMiddleware(http.HandlerFunc(s.handleListTransactions)))
	mux.Handle("/sep24/fee", authMiddleware(http.HandlerFunc(s.handleGetFee)))

	mux.HandleFunc("/sep24/interactive/deposit", s.renderDeposit)
	mux.HandleFunc("/sep24/interactive/withdraw", s.renderWithdraw)
	mux.HandleFunc("/sep24/interactive/status", s.renderStatus)
	mux.HandleFunc("/sep24/static/styles.css", s.renderStyle)
}

func accountFromRequest(r *http.Request) string {
	return middleware.AccountFromContext(r.Context())
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func parseTemplate(path string, fallback string) (*template.Template, error) {
	if _, err := os.Stat(path); err == nil {
		return template.ParseFiles(path)
	}
	return template.New("fallback").Parse(fallback)
}

func (s *Service) renderDeposit(w http.ResponseWriter, r *http.Request) {
	tpl, err := parseTemplate(filepath.FromSlash("sep24/interactive/templates/deposit.html"), "<html><body><h1>Deposit</h1></body></html>")
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("template error: %v", err))
		return
	}
	_ = tpl.Execute(w, map[string]string{"Title": "Deposit"})
}

func (s *Service) renderWithdraw(w http.ResponseWriter, r *http.Request) {
	tpl, err := parseTemplate(filepath.FromSlash("sep24/interactive/templates/withdraw.html"), "<html><body><h1>Withdraw</h1></body></html>")
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("template error: %v", err))
		return
	}
	_ = tpl.Execute(w, map[string]string{"Title": "Withdraw"})
}

func (s *Service) renderStatus(w http.ResponseWriter, r *http.Request) {
	tpl, err := parseTemplate(filepath.FromSlash("sep24/interactive/templates/status.html"), "<html><body><h1>Status</h1></body></html>")
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("template error: %v", err))
		return
	}
	_ = tpl.Execute(w, map[string]string{"Title": "Status", "ID": r.URL.Query().Get("id")})
}

func (s *Service) renderStyle(w http.ResponseWriter, r *http.Request) {
	path := filepath.FromSlash("sep24/interactive/static/styles.css")
	if _, err := os.Stat(path); err == nil {
		http.ServeFile(w, r, path)
		return
	}
	w.Header().Set("Content-Type", "text/css")
	_, _ = w.Write([]byte("body { font-family: sans-serif; margin: 2rem; }"))
}
