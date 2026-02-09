package sep10

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
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
	if homeDomain == "" {
		homeDomain = s.HomeDomain
	}
	return BuildChallenge(BuildParams{
		ServerAccount:     s.ServerAccount,
		ServerSigningKey:  s.ServerSigningKey,
		ClientAccount:     account,
		HomeDomain:        homeDomain,
		WebAuthDomain:     s.WebAuthDomain,
		ClientDomain:      clientDomain,
		NetworkPassphrase: s.NetworkPassphrase,
		Memo:              memo,
		TTL:               s.ChallengeTTL,
	})
}

func (s *Service) VerifyAndIssueToken(encodedChallenge string) (string, error) {
	result, err := VerifyChallenge(VerifyParams{
		EncodedChallenge: encodedChallenge,
		ServerAccount:    s.ServerAccount,
		ServerSigningKey: s.ServerSigningKey,
		RequireClientSig: true,
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json request")
		return
	}
	if strings.TrimSpace(req.Transaction) == "" {
		writeError(w, http.StatusBadRequest, "missing transaction")
		return
	}

	token, err := service.VerifyAndIssueToken(req.Transaction)
	if err != nil {
		writeError(w, http.StatusUnauthorized, fmt.Sprintf("challenge verification failed: %v", err))
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
