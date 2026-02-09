package sep10

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type Operation struct {
	Type   string `json:"type"`
	Source string `json:"source"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

type Challenge struct {
	SourceAccount     string      `json:"source_account"`
	Sequence          int64       `json:"sequence"`
	NetworkPassphrase string      `json:"network_passphrase"`
	HomeDomain        string      `json:"home_domain"`
	WebAuthDomain     string      `json:"web_auth_domain"`
	ClientDomain      string      `json:"client_domain,omitempty"`
	ClientAccount     string      `json:"client_account"`
	Memo              string      `json:"memo,omitempty"`
	IssuedAt          time.Time   `json:"issued_at"`
	ExpiresAt         time.Time   `json:"expires_at"`
	Operations        []Operation `json:"operations"`
}

type ChallengeEnvelope struct {
	Challenge  Challenge         `json:"challenge"`
	Signatures map[string]string `json:"signatures"`
}

type BuildParams struct {
	ServerAccount     string
	ServerSigningKey  string
	ClientAccount     string
	HomeDomain        string
	WebAuthDomain     string
	ClientDomain      string
	NetworkPassphrase string
	Memo              string
	Now               time.Time
	TTL               time.Duration
}

func BuildChallenge(params BuildParams) (string, error) {
	if params.ClientAccount == "" {
		return "", fmt.Errorf("client account is required")
	}
	if params.ServerAccount == "" {
		return "", fmt.Errorf("server account is required")
	}
	if params.HomeDomain == "" {
		return "", fmt.Errorf("home domain is required")
	}
	if params.WebAuthDomain == "" {
		params.WebAuthDomain = params.HomeDomain
	}
	if params.Now.IsZero() {
		params.Now = time.Now().UTC()
	}
	if params.TTL <= 0 {
		params.TTL = 5 * time.Minute
	}

	nonce := make([]byte, 64)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}

	challenge := Challenge{
		SourceAccount:     params.ServerAccount,
		Sequence:          0,
		NetworkPassphrase: params.NetworkPassphrase,
		HomeDomain:        params.HomeDomain,
		WebAuthDomain:     params.WebAuthDomain,
		ClientDomain:      params.ClientDomain,
		ClientAccount:     params.ClientAccount,
		Memo:              params.Memo,
		IssuedAt:          params.Now,
		ExpiresAt:         params.Now.Add(params.TTL),
		Operations: []Operation{
			{
				Type:   "manage_data",
				Source: params.ClientAccount,
				Name:   params.HomeDomain + " auth",
				Value:  base64.StdEncoding.EncodeToString(nonce),
			},
			{
				Type:   "manage_data",
				Source: params.ServerAccount,
				Name:   "web_auth_domain",
				Value:  params.WebAuthDomain,
			},
		},
	}

	payload, err := canonicalChallengeBytes(challenge)
	if err != nil {
		return "", err
	}

	envelope := ChallengeEnvelope{
		Challenge: challenge,
		Signatures: map[string]string{
			"server": sign(payload, params.ServerSigningKey),
		},
	}

	encoded, err := encodeEnvelope(envelope)
	if err != nil {
		return "", err
	}

	return encoded, nil
}

func AddClientSignature(encodedChallenge string, clientSigningSecret string) (string, error) {
	envelope, err := decodeEnvelope(encodedChallenge)
	if err != nil {
		return "", err
	}
	payload, err := canonicalChallengeBytes(envelope.Challenge)
	if err != nil {
		return "", err
	}
	if envelope.Signatures == nil {
		envelope.Signatures = map[string]string{}
	}
	envelope.Signatures["client"] = sign(payload, clientSigningSecret)
	return encodeEnvelope(envelope)
}

func canonicalChallengeBytes(challenge Challenge) ([]byte, error) {
	payload, err := json.Marshal(challenge)
	if err != nil {
		return nil, fmt.Errorf("marshal challenge: %w", err)
	}
	return payload, nil
}

func encodeEnvelope(envelope ChallengeEnvelope) (string, error) {
	raw, err := json.Marshal(envelope)
	if err != nil {
		return "", fmt.Errorf("marshal envelope: %w", err)
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}

func decodeEnvelope(encoded string) (ChallengeEnvelope, error) {
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return ChallengeEnvelope{}, fmt.Errorf("decode challenge envelope: %w", err)
	}
	var envelope ChallengeEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return ChallengeEnvelope{}, fmt.Errorf("unmarshal challenge envelope: %w", err)
	}
	return envelope, nil
}

func sign(payload []byte, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write(payload)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
