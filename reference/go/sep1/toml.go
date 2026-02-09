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

	for _, asset := range cfg.Assets {
		b.WriteString("\n[[CURRENCIES]]\n")
		b.WriteString(fmt.Sprintf("code=\"%s\"\n", asset.Code))
		b.WriteString("status=\"live\"\n")
		b.WriteString("desc=\"Reference anchor asset\"\n")
		b.WriteString("is_asset_anchored=true\n")
		b.WriteString("anchor_asset_type=\"fiat\"\n")
		b.WriteString(fmt.Sprintf("anchor_asset=\"%s\"\n", asset.Code))
	}

	return b.String()
}
