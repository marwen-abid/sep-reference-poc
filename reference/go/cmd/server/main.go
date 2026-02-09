package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stellar/sep-reference/reference/go/internal/config"
	"github.com/stellar/sep-reference/reference/go/internal/db"
	"github.com/stellar/sep-reference/reference/go/internal/middleware"
	"github.com/stellar/sep-reference/reference/go/sep1"
	"github.com/stellar/sep-reference/reference/go/sep10"
	"github.com/stellar/sep-reference/reference/go/sep24"
)

func main() {
	cfg := config.Load()

	txStore := db.NewMemoryTransactionStore()
	customerStore := db.NewMemoryCustomerStore()

	authService := sep10.NewService(
		cfg.ServerAccount,
		cfg.SigningKey,
		cfg.JWTSecret,
		cfg.HomeDomain,
		cfg.WebAuthDomain,
		cfg.NetworkPassphrase,
		cfg.ChallengeTTL,
		cfg.TokenTTL,
	)

	sep24Service := sep24.NewService(cfg, txStore, customerStore)

	mux := http.NewServeMux()
	mux.Handle("/auth", sep10.NewHTTPHandler(authService))
	mux.HandleFunc("/.well-known/stellar.toml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(sep1.RenderStellarTOML(cfg)))
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	sep24Service.RegisterRoutes(mux, middleware.SEP10Auth(cfg.JWTSecret))

	log.Printf("SEP Reference server starting")
	log.Printf("SEP-1:  http://%s/.well-known/stellar.toml", cfg.HomeDomain)
	log.Printf("SEP-10: http://%s/auth", cfg.HomeDomain)
	log.Printf("SEP-24: http://%s/sep24", cfg.HomeDomain)
	log.Printf("Ready to accept connections")

	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatal(fmt.Errorf("server exited: %w", err))
	}
}
