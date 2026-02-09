package sep10

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/stellar/go/xdr"
)

type Service struct {
	ServerAccount     string
	ServerSigningKey  string
	JWTSecret         string
	HomeDomain        string
	WebAuthDomain     string
	NetworkPassphrase string
	ChallengeTTL      time.Duration
	TokenTTL          time.Duration
	AccountSigners    AccountSignerLoader
}

type challengeResponse struct {
	Transaction       string `json:"transaction"`
	NetworkPassphrase string `json:"network_passphrase"`
}

type verifyRequest struct {
	Transaction string `json:"transaction"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewService(serverAccount, serverSigningKey, jwtSecret, homeDomain, webAuthDomain, networkPassphrase string, challengeTTL, tokenTTL time.Duration) *Service {
	return &Service{
		ServerAccount:     serverAccount,
		ServerSigningKey:  serverSigningKey,
		JWTSecret:         jwtSecret,
		HomeDomain:        homeDomain,
		WebAuthDomain:     webAuthDomain,
		NetworkPassphrase: networkPassphrase,
		ChallengeTTL:      challengeTTL,
		TokenTTL:          tokenTTL,
	}
}

func (s *Service) BuildChallenge(account, clientDomain, homeDomain, memo string) (string, error) {
	_ = clientDomain
	if homeDomain == "" {
		homeDomain = s.HomeDomain
	}
	return BuildChallenge(BuildParams{
		ServerSigningKey:  s.ServerSigningKey,
		ClientAccount:     account,
		HomeDomain:        homeDomain,
		WebAuthDomain:     s.WebAuthDomain,
		NetworkPassphrase: s.NetworkPassphrase,
		Memo:              memo,
		TTL:               s.ChallengeTTL,
	})
}

func (s *Service) VerifyAndIssueToken(encodedChallenge string) (string, error) {
	result, err := VerifyChallenge(VerifyParams{
		EncodedChallenge:  encodedChallenge,
		ServerAccount:     s.ServerAccount,
		NetworkPassphrase: s.NetworkPassphrase,
		WebAuthDomain:     s.WebAuthDomain,
		HomeDomains:       []string{s.HomeDomain},
		RequireClientSig:  true,
		AccountSigners:    s.AccountSigners,
	})
	if err != nil {
		return "", err
	}
	return IssueToken(result.ClientAccount, s.HomeDomain, result.ClientDomain, result.HomeDomain, s.JWTSecret, time.Now().UTC(), s.TokenTTL)
}

func NewHTTPHandler(service *Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetAuth(w, r, service)
		case http.MethodPost:
			handlePostAuth(w, r, service)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})
	return mux
}

func handleGetAuth(w http.ResponseWriter, r *http.Request, service *Service) {
	account := strings.TrimSpace(r.URL.Query().Get("account"))
	if account == "" {
		writeError(w, http.StatusBadRequest, "missing account")
		return
	}
	if !isValidStellarAddress(account) {
		writeError(w, http.StatusBadRequest, "invalid account")
		return
	}
	clientDomain := strings.TrimSpace(r.URL.Query().Get("client_domain"))
	homeDomain := strings.TrimSpace(r.URL.Query().Get("home_domain"))
	memo := strings.TrimSpace(r.URL.Query().Get("memo"))

	challenge, err := service.BuildChallenge(account, clientDomain, homeDomain, memo)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("build challenge: %v", err))
		return
	}
	writeJSON(w, http.StatusOK, challengeResponse{
		Transaction:       challenge,
		NetworkPassphrase: service.NetworkPassphrase,
	})
}

func handlePostAuth(w http.ResponseWriter, r *http.Request, service *Service) {
	defer r.Body.Close()
	var req verifyRequest
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		writeError(w, http.StatusBadRequest, "missing transaction")
		return
	}

	if err := json.Unmarshal(raw, &req); err != nil {
		form, formErr := url.ParseQuery(string(raw))
		if formErr != nil {
			writeError(w, http.StatusBadRequest, "invalid request payload")
			return
		}
		req.Transaction = strings.TrimSpace(form.Get("transaction"))
	}
	if strings.TrimSpace(req.Transaction) == "" {
		writeError(w, http.StatusBadRequest, "missing transaction")
		return
	}

	token, err := service.VerifyAndIssueToken(req.Transaction)
	if err != nil {
		writeError(w, statusCodeForVerifyError(err), fmt.Sprintf("challenge verification failed: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{Token: token})
}

func writeJSON(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, errorResponse{Error: message})
}

func isValidStellarAddress(raw string) bool {
	if _, err := xdr.AddressToAccountId(raw); err == nil {
		return true
	}
	if _, err := xdr.AddressToMuxedAccount(raw); err == nil {
		return true
	}
	return false
}
