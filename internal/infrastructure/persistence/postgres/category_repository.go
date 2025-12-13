package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/taketosaeki/donelog/internal/app/donelog/command"
	"github.com/taketosaeki/donelog/internal/domain/donelog"
)

// CategoryRepository provides Category lookup backed by Postgres.
type CategoryRepository struct {
	pool *pgxpool.Pool
}

// NewCategoryRepository constructs a CategoryRepository.
func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

// FindActiveByID returns a Category if found, or nil if absent.
func (r *CategoryRepository) FindActiveByID(ctx context.Context, id donelog.CategoryID) (*command.Category, error) {
	const query = `
SELECT id, active
FROM categories
WHERE id = $1;
`

	row := r.pool.QueryRow(ctx, query, id.String())

	var (
		rawID  string
		active bool
	)

	if err := row.Scan(&rawID, &active); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	catID, err := donelog.NewCategoryID(rawID)
	if err != nil {
		return nil, err
	}

	return &command.Category{
		ID:     catID,
		Active: active,
	}, nil
}

var _ command.CategoryRepository = (*CategoryRepository)(nil)
