package sep24

import "net/http"

type assetInfo struct {
	Enabled    bool   `json:"enabled"`
	FeeFixed   string `json:"fee_fixed"`
	FeePercent string `json:"fee_percent"`
}

func (s *Service) handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	deposit := map[string]assetInfo{}
	withdraw := map[string]assetInfo{}
	for _, asset := range s.Config.Assets {
		info := assetInfo{
			Enabled:    asset.Enabled,
			FeeFixed:   formatAmount(asset.FeeFixed),
			FeePercent: formatAmount(asset.FeePercent),
		}
		deposit[asset.Code] = info
		withdraw[asset.Code] = info
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"deposit":  deposit,
		"withdraw": withdraw,
	})
}
