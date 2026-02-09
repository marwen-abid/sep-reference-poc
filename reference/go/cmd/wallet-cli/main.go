package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/stellar/sep-reference/reference/go/sep10"
)

type challengeResponse struct {
	Transaction       string `json:"transaction"`
	NetworkPassphrase string `json:"network_passphrase"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func main() {
	baseURL := flag.String("base-url", "http://localhost:8080", "server base URL")
	account := flag.String("account", "GCTSEPLDTRBIDEZ2WWOECTX5QPRTKD4GRDCRQBHWRDBTJROWJG2G6MGV", "client stellar account")
	clientSecret := flag.String("client-secret", "SCFDN4SWA4VR2Z2FDMGSQSTIYKNAL7LLWD6LCBZ7OTZ4LORMHXY2HUT4", "client signing secret")
	flag.Parse()

	tx, networkPassphrase, err := getChallenge(*baseURL, *account)
	if err != nil {
		panic(err)
	}

	signed, err := sep10.AddClientSignatureWithNetworkPassphrase(tx, *clientSecret, networkPassphrase)
	if err != nil {
		panic(err)
	}

	token, err := postChallenge(*baseURL, signed)
	if err != nil {
		panic(err)
	}

	fmt.Println(token)
}

func getChallenge(baseURL string, account string) (string, string, error) {
	u, _ := url.Parse(baseURL + "/auth")
	q := u.Query()
	q.Set("account", account)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("get challenge failed: %s", string(raw))
	}

	var out challengeResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", "", err
	}
	return out.Transaction, out.NetworkPassphrase, nil
}

func postChallenge(baseURL string, challenge string) (string, error) {
	payload, _ := json.Marshal(map[string]string{"transaction": challenge})
	resp, err := http.Post(baseURL+"/auth", "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("post challenge failed: %s", string(raw))
	}

	var out tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Token, nil
}
