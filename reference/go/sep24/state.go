package sep24

import "fmt"

const (
	StatusIncomplete               = "incomplete"
	StatusPendingUserTransferStart = "pending_user_transfer_start"
	StatusPendingAnchor            = "pending_anchor"
	StatusPendingStellar           = "pending_stellar"
	StatusCompleted                = "completed"
	StatusError                    = "error"
	StatusExpired                  = "expired"
)

var allowedTransitions = map[string]map[string]bool{
	StatusIncomplete: {
		StatusPendingUserTransferStart: true,
		StatusExpired:                  true,
		StatusError:                    true,
	},
	StatusPendingUserTransferStart: {
		StatusPendingAnchor: true,
		StatusExpired:       true,
		StatusError:         true,
	},
	StatusPendingAnchor: {
		StatusPendingStellar: true,
		StatusError:          true,
	},
	StatusPendingStellar: {
		StatusCompleted: true,
		StatusError:     true,
	},
}

func ValidateTransition(from, to string) error {
	if from == to {
		return nil
	}
	if next, ok := allowedTransitions[from]; ok && next[to] {
		return nil
	}
	return fmt.Errorf("invalid transition from %s to %s", from, to)
}
