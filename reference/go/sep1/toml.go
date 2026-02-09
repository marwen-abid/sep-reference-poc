package sep1

import (
	"fmt"
	"strings"

	"github.com/stellar/sep-reference/reference/go/internal/config"
)

func RenderStellarTOML(cfg config.Config) string {
	var b strings.Builder
	b.WriteString("VERSION=\"2.7.0\"\n")
	b.WriteString(fmt.Sprintf("NETWORK_PASSPHRASE=\"%s\"\n", cfg.NetworkPassphrase))
	b.WriteString(fmt.Sprintf("SIGNING_KEY=\"%s\"\n", cfg.ServerAccount))
	b.WriteString(fmt.Sprintf("WEB_AUTH_ENDPOINT=\"http://%s/auth\"\n", cfg.HomeDomain))
	b.WriteString(fmt.Sprintf("TRANSFER_SERVER=\"%s\"\n", cfg.TransferServer))
	b.WriteString(fmt.Sprintf("TRANSFER_SERVER_SEP0024=\"%s\"\n", cfg.TransferServerSep24))
	if cfg.QuoteServer != "" {
		b.WriteString(fmt.Sprintf("ANCHOR_QUOTE_SERVER=\"%s\"\n", cfg.QuoteServer))
	}
	return b.String()
}
