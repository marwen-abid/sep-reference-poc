package sep10

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Claims struct {
	Subject      string `json:"sub"`
	Issuer       string `json:"iss"`
	IssuedAt     int64  `json:"iat"`
	ExpiresAt    int64  `json:"exp"`
	JWTID        string `json:"jti"`
	ClientDomain string `json:"client_domain,omitempty"`
	HomeDomain   string `json:"home_domain,omitempty"`
}

func IssueToken(sub, iss, clientDomain, homeDomain, secret string, now time.Time, ttl time.Duration) (string, error) {
	if sub == "" {
		return "", fmt.Errorf("subject is required")
	}
	if iss == "" {
		return "", fmt.Errorf("issuer is required")
	}
	if secret == "" {
		return "", fmt.Errorf("secret is required")
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}

	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	jti, err := randomTokenID()
	if err != nil {
		return "", err
	}
	claims := Claims{
		Subject:      sub,
		Issuer:       iss,
		IssuedAt:     now.Unix(),
		ExpiresAt:    now.Add(ttl).Unix(),
		JWTID:        jti,
		ClientDomain: clientDomain,
		HomeDomain:   homeDomain,
	}

	headerJSON, _ := json.Marshal(header)
	claimsJSON, _ := json.Marshal(claims)
	headerPart := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsPart := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := headerPart + "." + claimsPart

	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(signingInput))
	sigPart := base64.RawURLEncoding.EncodeToString(sig.Sum(nil))

	return signingInput + "." + sigPart, nil
}

func VerifyToken(token, secret string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, fmt.Errorf("invalid token format")
	}

	signingInput := parts[0] + "." + parts[1]
	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(signingInput))
	expected := base64.RawURLEncoding.EncodeToString(sig.Sum(nil))
	if !hmac.Equal([]byte(parts[2]), []byte(expected)) {
		return Claims{}, fmt.Errorf("invalid token signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, fmt.Errorf("invalid token payload")
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, fmt.Errorf("invalid token claims")
	}
	if time.Now().UTC().Unix() > claims.ExpiresAt {
		return Claims{}, fmt.Errorf("token expired")
	}
	return claims, nil
}

func randomTokenID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate jti: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
