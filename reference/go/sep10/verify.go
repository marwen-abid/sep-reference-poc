package sep10

import (
	"errors"
	"fmt"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

type VerifyParams struct {
	EncodedChallenge  string
	ServerAccount     string
	NetworkPassphrase string
	WebAuthDomain     string
	HomeDomains       []string
	RequireClientSig  bool
	AccountSigners    AccountSignerLoader
}

type VerifyResult struct {
	ClientAccount string
	ClientDomain  string
	HomeDomain    string
}

type AccountSignerLoader func(accountID string) (txnbuild.SignerSummary, txnbuild.Threshold, bool, error)

type verifyError struct {
	Status int
	Err    error
}

func (e *verifyError) Error() string {
	return e.Err.Error()
}

func (e *verifyError) Unwrap() error {
	return e.Err
}

func invalidChallenge(err error) error {
	return &verifyError{Status: 400, Err: err}
}

func unauthorizedChallenge(err error) error {
	return &verifyError{Status: 400, Err: err}
}

func statusCodeForVerifyError(err error) int {
	var vErr *verifyError
	if errors.As(err, &vErr) {
		return vErr.Status
	}
	return 400
}

func VerifyChallenge(params VerifyParams) (VerifyResult, error) {
	if params.NetworkPassphrase == "" {
		params.NetworkPassphrase = DefaultNetworkPassphrase
	}
	if len(params.HomeDomains) == 0 {
		return VerifyResult{}, invalidChallenge(fmt.Errorf("at least one home domain is required"))
	}

	_, clientAccount, matchedHomeDomain, _, err := txnbuild.ReadChallengeTx(
		params.EncodedChallenge,
		params.ServerAccount,
		params.NetworkPassphrase,
		params.WebAuthDomain,
		params.HomeDomains,
	)
	if err != nil {
		return VerifyResult{}, invalidChallenge(fmt.Errorf("invalid challenge transaction: %w", err))
	}

	result := VerifyResult{
		ClientAccount: clientAccount,
		HomeDomain:    matchedHomeDomain,
	}
	if !params.RequireClientSig {
		return result, nil
	}

	if err := verifyClientSignatures(params, clientAccount); err != nil {
		return VerifyResult{}, err
	}
	return result, nil
}

func verifyClientSignatures(params VerifyParams, clientAccount string) error {
	accountSigners := params.AccountSigners
	if accountSigners == nil {
		accountSigners = buildDefaultAccountSignerLoader(params.NetworkPassphrase)
	}

	signerSummary, threshold, exists, err := accountSigners(clientAccount)
	if err != nil {
		return unauthorizedChallenge(fmt.Errorf("load account signers: %w", err))
	}

	if !exists {
		_, verifyErr := txnbuild.VerifyChallengeTxSigners(
			params.EncodedChallenge,
			params.ServerAccount,
			params.NetworkPassphrase,
			params.WebAuthDomain,
			params.HomeDomains,
			clientAccount,
		)
		if verifyErr != nil {
			return unauthorizedChallenge(fmt.Errorf("challenge verification failed: %w", verifyErr))
		}
		return nil
	}

	_, verifyErr := txnbuild.VerifyChallengeTxThreshold(
		params.EncodedChallenge,
		params.ServerAccount,
		params.NetworkPassphrase,
		params.WebAuthDomain,
		params.HomeDomains,
		threshold,
		signerSummary,
	)
	if verifyErr != nil {
		return unauthorizedChallenge(fmt.Errorf("challenge verification failed: %w", verifyErr))
	}
	return nil
}

func buildDefaultAccountSignerLoader(networkPassphrase string) AccountSignerLoader {
	hClient := horizonForNetwork(networkPassphrase)
	return func(accountID string) (txnbuild.SignerSummary, txnbuild.Threshold, bool, error) {
		account, err := hClient.AccountDetail(horizonclient.AccountRequest{AccountID: accountID})
		if err != nil {
			if horizonclient.IsNotFoundError(err) {
				return nil, 0, false, nil
			}
			return nil, 0, false, err
		}

		signerSummary := txnbuild.SignerSummary(account.SignerSummary())
		threshold := txnbuild.Threshold(account.Thresholds.MedThreshold)
		if threshold < 1 {
			threshold = 1
		}

		return signerSummary, threshold, true, nil
	}
}

func horizonForNetwork(networkPassphrase string) *horizonclient.Client {
	switch networkPassphrase {
	case network.TestNetworkPassphrase:
		return horizonclient.DefaultTestNetClient
	case network.PublicNetworkPassphrase:
		return horizonclient.DefaultPublicNetClient
	default:
		return horizonclient.DefaultTestNetClient
	}
}
