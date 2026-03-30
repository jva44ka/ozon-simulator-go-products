package reservation

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *Service) Release(ctx context.Context, ids []int64) error {
	return s.db.InTransaction(ctx, func(tx pgx.Tx) error {
		if err := s.db.Reservations().WithTx(tx).DeleteByIds(ctx, ids); err != nil {
			return fmt.Errorf("Release: %w", err)
		}
		return nil
	})
}
