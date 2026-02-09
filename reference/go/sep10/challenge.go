package sep10

import (
	"fmt"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

const DefaultNetworkPassphrase = "Test SDF Network ; September 2015"

type BuildParams struct {
	ServerSigningKey  string
	ClientAccount     string
	HomeDomain        string
	WebAuthDomain     string
	NetworkPassphrase string
	Memo              string
	TTL               time.Duration
}

func BuildChallenge(params BuildParams) (string, error) {
	if params.ClientAccount == "" {
		return "", fmt.Errorf("client account is required")
	}
	if params.ServerSigningKey == "" {
		return "", fmt.Errorf("server signing key is required")
	}
	if params.HomeDomain == "" {
		return "", fmt.Errorf("home domain is required")
	}
	if params.WebAuthDomain == "" {
		params.WebAuthDomain = params.HomeDomain
	}
	if params.NetworkPassphrase == "" {
		params.NetworkPassphrase = DefaultNetworkPassphrase
	}
	if params.TTL <= 0 {
		params.TTL = 5 * time.Minute
	}

	var memo *txnbuild.MemoID
	if params.Memo != "" {
		return "", fmt.Errorf("memo is not supported by this server")
	}

	tx, err := txnbuild.BuildChallengeTx(
		params.ServerSigningKey,
		params.ClientAccount,
		params.WebAuthDomain,
		params.HomeDomain,
		params.NetworkPassphrase,
		params.TTL,
		memo,
	)
	if err != nil {
		return "", fmt.Errorf("build challenge tx: %w", err)
	}

	encoded, err := tx.Base64()
	if err != nil {
		return "", fmt.Errorf("encode challenge tx: %w", err)
	}
	return encoded, nil
}

func AddClientSignature(encodedChallenge, clientSigningSecret string) (string, error) {
	return AddClientSignatureWithNetworkPassphrase(encodedChallenge, clientSigningSecret, DefaultNetworkPassphrase)
}

func AddClientSignatureWithNetworkPassphrase(encodedChallenge, clientSigningSecret, networkPassphrase string) (string, error) {
	if networkPassphrase == "" {
		networkPassphrase = DefaultNetworkPassphrase
	}

	parsed, err := txnbuild.TransactionFromXDR(encodedChallenge)
	if err != nil {
		return "", fmt.Errorf("parse challenge tx: %w", err)
	}
	tx, ok := parsed.Transaction()
	if !ok {
		return "", fmt.Errorf("challenge must be a non-fee-bump transaction")
	}

	clientKP, err := keypair.ParseFull(clientSigningSecret)
	if err != nil {
		return "", fmt.Errorf("parse client signing key: %w", err)
	}

	signed, err := tx.Sign(networkPassphrase, clientKP)
	if err != nil {
		return "", fmt.Errorf("sign challenge tx: %w", err)
	}
	encoded, err := signed.Base64()
	if err != nil {
		return "", fmt.Errorf("encode signed challenge tx: %w", err)
	}
	return encoded, nil
}
