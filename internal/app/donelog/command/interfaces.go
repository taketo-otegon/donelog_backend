package command

import (
	"context"

	"github.com/taketosaeki/donelog/internal/domain/donelog"
)

// DoneLogRepository is the command-side abstraction for persistence.
type DoneLogRepository interface {
	Save(ctx context.Context, log *donelog.DoneLog) error
	FindByID(ctx context.Context, id donelog.DoneLogID) (*donelog.RawDoneLog, error)
	Delete(ctx context.Context, id donelog.DoneLogID) error
}

// TrackRepository provides access to Track aggregates.
type TrackRepository interface {
	FindActiveByID(ctx context.Context, id donelog.TrackID) (*Track, error)
}

// CategoryRepository provides access to Category aggregates.
type CategoryRepository interface {
	FindActiveByID(ctx context.Context, id donelog.CategoryID) (*Category, error)
}

// IDGenerator creates unique DoneLogID values.
type IDGenerator interface {
	NewDoneLogID(ctx context.Context) (donelog.DoneLogID, error)
}

// Clock provides current time for OccurredOn defaults.
type Clock interface {
	Now() donelog.OccurredOn
}
