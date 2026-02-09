package sep24

import "net/http"

type assetInfo struct {
	Enabled    bool    `json:"enabled"`
	FeeFixed   float64 `json:"fee_fixed"`
	FeePercent float64 `json:"fee_percent"`
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
			FeeFixed:   asset.FeeFixed,
			FeePercent: asset.FeePercent,
		}
		deposit[asset.Code] = info
		withdraw[asset.Code] = info
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"deposit":  deposit,
		"withdraw": withdraw,
		"fee": map[string]bool{
			"enabled": true,
		},
	})
}
