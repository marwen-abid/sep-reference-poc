package sep10

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"time"
)

type VerifyParams struct {
	EncodedChallenge string
	ServerAccount    string
	ServerSigningKey string
	Now              time.Time
	RequireClientSig bool
}

type VerifyResult struct {
	ClientAccount string
	ClientDomain  string
	HomeDomain    string
}

func VerifyChallenge(params VerifyParams) (VerifyResult, error) {
	envelope, err := decodeEnvelope(params.EncodedChallenge)
	if err != nil {
		return VerifyResult{}, err
	}
	challenge := envelope.Challenge

	if params.Now.IsZero() {
		params.Now = time.Now().UTC()
	}

	if challenge.SourceAccount != params.ServerAccount {
		return VerifyResult{}, fmt.Errorf("invalid source account")
	}
	if challenge.Sequence != 0 {
		return VerifyResult{}, fmt.Errorf("invalid sequence")
	}
	if params.Now.After(challenge.ExpiresAt) {
		return VerifyResult{}, fmt.Errorf("challenge expired")
	}
	if params.Now.Before(challenge.IssuedAt.Add(-1 * time.Minute)) {
		return VerifyResult{}, fmt.Errorf("challenge issued-at is in the future")
	}
	if len(challenge.Operations) == 0 {
		return VerifyResult{}, fmt.Errorf("challenge has no operations")
	}

	first := challenge.Operations[0]
	if first.Type != "manage_data" {
		return VerifyResult{}, fmt.Errorf("invalid first operation type")
	}
	if first.Source != challenge.ClientAccount {
		return VerifyResult{}, fmt.Errorf("first operation source mismatch")
	}
	if first.Name != challenge.HomeDomain+" auth" {
		return VerifyResult{}, fmt.Errorf("first operation name mismatch")
	}
	nonce, err := base64.StdEncoding.DecodeString(first.Value)
	if err != nil {
		return VerifyResult{}, fmt.Errorf("invalid nonce encoding")
	}
	if len(nonce) != 64 {
		return VerifyResult{}, fmt.Errorf("invalid nonce length")
	}

	payload, err := canonicalChallengeBytes(challenge)
	if err != nil {
		return VerifyResult{}, err
	}
	expectedServerSig := sign(payload, params.ServerSigningKey)
	serverSig, ok := envelope.Signatures["server"]
	if !ok {
		return VerifyResult{}, fmt.Errorf("missing server signature")
	}
	if subtle.ConstantTimeCompare([]byte(serverSig), []byte(expectedServerSig)) != 1 {
		return VerifyResult{}, fmt.Errorf("invalid server signature")
	}

	if params.RequireClientSig {
		clientSig, ok := envelope.Signatures["client"]
		if !ok || clientSig == "" {
			return VerifyResult{}, fmt.Errorf("missing client signature")
		}
	}

	return VerifyResult{
		ClientAccount: challenge.ClientAccount,
		ClientDomain:  challenge.ClientDomain,
		HomeDomain:    challenge.HomeDomain,
	}, nil
}
