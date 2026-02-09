package config

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"time"
)

type Asset struct {
	Code       string  `json:"asset_code"`
	Enabled    bool    `json:"enabled"`
	FeeFixed   float64 `json:"fee_fixed"`
	FeePercent float64 `json:"fee_percent"`
}

type Config struct {
	Addr                string
	HomeDomain          string
	WebAuthDomain       string
	NetworkPassphrase   string
	SigningKey          string
	ServerAccount       string
	JWTSecret           string
	ChallengeTTL        time.Duration
	TokenTTL            time.Duration
	TransferServer      string
	TransferServerSep24 string
	QuoteServer         string
	Assets              []Asset
}

func Load() Config {
	homeDomain := getenv("HOME_DOMAIN", "localhost:8080")
	transferServer := getenv("TRANSFER_SERVER", "http://localhost:8080/sep24")

	cfg := Config{
		Addr:                getenv("ADDR", ":8080"),
		HomeDomain:          homeDomain,
		WebAuthDomain:       getenv("WEB_AUTH_DOMAIN", homeDomain),
		NetworkPassphrase:   getenv("NETWORK_PASSPHRASE", "Test SDF Network ; September 2015"),
		SigningKey:          getenv("SIGNING_KEY", "dev-signing-key"),
		JWTSecret:           getenv("JWT_SECRET", "dev-jwt-secret"),
		ChallengeTTL:        parseDuration(getenv("CHALLENGE_TTL", "5m"), 5*time.Minute),
		TokenTTL:            parseDuration(getenv("TOKEN_TTL", "15m"), 15*time.Minute),
		TransferServer:      transferServer,
		TransferServerSep24: getenv("TRANSFER_SERVER_SEP0024", transferServer),
		QuoteServer:         getenv("QUOTE_SERVER", ""),
		Assets:              parseAssets(getenv("ASSETS", "USDC")),
	}

	cfg.ServerAccount = getenv("SERVER_ACCOUNT", derivePseudoAccount(cfg.SigningKey))
	return cfg
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(raw string, fallback time.Duration) time.Duration {
	d, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return d
}

func parseAssets(raw string) []Asset {
	parts := strings.Split(raw, ",")
	assets := make([]Asset, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		code, fixed, percent := parseAssetItem(item)
		assets = append(assets, Asset{Code: code, Enabled: true, FeeFixed: fixed, FeePercent: percent})
	}
	if len(assets) == 0 {
		assets = append(assets, Asset{Code: "USDC", Enabled: true, FeeFixed: 1.0, FeePercent: 0.0})
	}
	return assets
}

func parseAssetItem(item string) (string, float64, float64) {
	chunks := strings.Split(item, ":")
	code := strings.TrimSpace(chunks[0])
	fixed := 1.0
	percent := 0.0
	if len(chunks) > 1 {
		if v, err := strconv.ParseFloat(chunks[1], 64); err == nil {
			fixed = v
		}
	}
	if len(chunks) > 2 {
		if v, err := strconv.ParseFloat(chunks[2], 64); err == nil {
			percent = v
		}
	}
	return code, fixed, percent
}

func derivePseudoAccount(seed string) string {
	sum := sha256.Sum256([]byte(seed))
	encoded := strings.ToUpper(hex.EncodeToString(sum[:]))
	if len(encoded) < 55 {
		encoded = encoded + strings.Repeat("A", 55-len(encoded))
	}
	return "G" + encoded[:55]
}
