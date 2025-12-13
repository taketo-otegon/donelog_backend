package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/taketosaeki/donelog/internal/app/donelog/command"
	"github.com/taketosaeki/donelog/internal/domain/donelog"
)

// TrackRepository provides Track lookup backed by Postgres.
type TrackRepository struct {
	pool *pgxpool.Pool
}

// NewTrackRepository constructs a TrackRepository.
func NewTrackRepository(pool *pgxpool.Pool) *TrackRepository {
	return &TrackRepository{pool: pool}
}

// FindActiveByID returns a Track if found, or nil if absent.
func (r *TrackRepository) FindActiveByID(ctx context.Context, id donelog.TrackID) (*command.Track, error) {
	const query = `
SELECT id, default_category_id, active
FROM tracks
WHERE id = $1;
`

	row := r.pool.QueryRow(ctx, query, id.String())

	var (
		rawID     string
		defaultID *string
		active    bool
	)

	if err := row.Scan(&rawID, &defaultID, &active); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	trackID, err := donelog.NewTrackID(rawID)
	if err != nil {
		return nil, err
	}

	var defaultCategory *donelog.CategoryID
	if defaultID != nil {
		catID, err := donelog.NewCategoryID(*defaultID)
		if err != nil {
			return nil, err
		}
		defaultCategory = &catID
	}

	return &command.Track{
		ID:              trackID,
		DefaultCategory: defaultCategory,
		Active:          active,
	}, nil
}

var _ command.TrackRepository = (*TrackRepository)(nil)
