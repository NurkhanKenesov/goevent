
package invite

import (
	"context"

	"goevent/internal/logging"
)

func StartInviteWorker(ctx context.Context, jobs <-chan InviteJob) {
	logging.Info("invite.worker.start", nil)
	for {
		select {
		case <-ctx.Done():
			logging.Info("invite.worker.stop", nil)
			return
		case job := <-jobs:
			// attempt to include request_id if present on job (best-effort)
			fields := map[string]interface{}{"event_id": job.EventID, "email": job.Email, "user_id": job.UserID}
			if rid := jobRequestID(job); rid != "" {
				fields["request_id"] = rid
			}
			logging.Info("invite.worker.send", fields)
		}
	}
}

// helper to extract request id if job type has it (compile-time check expects field)
func jobRequestID(j InviteJob) string {
	// if InviteJob has RequestID field it will be accessible; otherwise return empty
	type withRID interface{ GetRequestID() string }
	if v, ok := any(j).(withRID); ok {
		return v.GetRequestID()
	}
	return ""
}

