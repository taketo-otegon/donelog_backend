package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/taketosaeki/donelog/internal/app/donelog/command"
	"github.com/taketosaeki/donelog/internal/domain/donelog"
)

// DoneLogRepository persists DoneLog aggregates in Postgres.
type DoneLogRepository struct {
	pool *pgxpool.Pool
}

// NewDoneLogRepository constructs a repository backed by pgxpool.
func NewDoneLogRepository(pool *pgxpool.Pool) *DoneLogRepository {
	return &DoneLogRepository{pool: pool}
}

// Save inserts or updates a DoneLog record.
func (r *DoneLogRepository) Save(ctx context.Context, log *donelog.DoneLog) error {
	const query = `
INSERT INTO donelogs (
    id, title, track_id, category_id, count, occurred_on, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, NOW(), NOW()
) ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    track_id = EXCLUDED.track_id,
    category_id = EXCLUDED.category_id,
    count = EXCLUDED.count,
    occurred_on = EXCLUDED.occurred_on,
    updated_at = NOW();
`

	_, err := r.pool.Exec(ctx, query,
		log.ID().String(),
		log.Title().String(),
		log.TrackID().String(),
		log.CategoryID().String(),
		log.Count().Int(),
		log.OccurredOn().Time(),
	)
	return err
}

// FindByID returns RawDoneLog for rehydration or nil if not found.
func (r *DoneLogRepository) FindByID(ctx context.Context, id donelog.DoneLogID) (*donelog.RawDoneLog, error) {
	const query = `
SELECT id, title, track_id, category_id, count, occurred_on
FROM donelogs
WHERE id = $1;
`

	row := r.pool.QueryRow(ctx, query, id.String())
	var raw donelog.RawDoneLog
	if err := row.Scan(&raw.ID, &raw.Title, &raw.TrackID, &raw.CategoryID, &raw.Count, &raw.OccurredOn); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &raw, nil
}

// Delete removes a DoneLog by ID. No error is returned if the row is absent.
func (r *DoneLogRepository) Delete(ctx context.Context, id donelog.DoneLogID) error {
	const query = `DELETE FROM donelogs WHERE id = $1;`
	_, err := r.pool.Exec(ctx, query, id.String())
	return err
}

var _ command.DoneLogRepository = (*DoneLogRepository)(nil)
