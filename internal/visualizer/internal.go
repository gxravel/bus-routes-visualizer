package visualizer

import "context"

func (r *Visualizer) IsHealthy(ctx context.Context) error {
	return r.db.StatusCheck(ctx)
}
