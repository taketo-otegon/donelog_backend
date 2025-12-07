package command

import (
	"context"
	"fmt"

	"github.com/taketosaeki/donelog/internal/domain/donelog"
)

// CreateDoneLogCommand holds the input data for creating a DONELOG.
type CreateDoneLogCommand struct {
	Title      string
	TrackID    string
	CategoryID string
	Count      int
	OccurredOn string
}

// Validate performs basic checks before constructing VO.
func (c CreateDoneLogCommand) Validate() error {
	if c.Title == "" {
		return fmt.Errorf("title is required")
	}
	if c.TrackID == "" {
		return fmt.Errorf("trackId is required")
	}
	if c.CategoryID == "" {
		return fmt.Errorf("categoryId is required")
	}
	if c.Count <= 0 {
		return fmt.Errorf("count must be > 0")
	}
	if c.OccurredOn == "" {
		return fmt.Errorf("occurredOn is required")
	}
	return nil
}

// CreateDoneLogHandler handles CreateDoneLogCommand.
type CreateDoneLogHandler struct {
	DoneLogs   DoneLogRepository
	Tracks     TrackRepository
	Categories CategoryRepository
	IDs        IDGenerator
}

// Handle executes the command and returns the new DoneLogID.
func (h CreateDoneLogHandler) Handle(ctx context.Context, cmd CreateDoneLogCommand) (donelog.DoneLogID, error) {
	if err := cmd.Validate(); err != nil {
		return donelog.DoneLogID{}, err
	}

	title, err := donelog.NewTitle(cmd.Title)
	if err != nil {
		return donelog.DoneLogID{}, err
	}
	trackID, err := donelog.NewTrackID(cmd.TrackID)
	if err != nil {
		return donelog.DoneLogID{}, err
	}
	categoryID, err := donelog.NewCategoryID(cmd.CategoryID)
	if err != nil {
		return donelog.DoneLogID{}, err
	}
	count, err := donelog.NewCount(cmd.Count)
	if err != nil {
		return donelog.DoneLogID{}, err
	}

	track, err := h.Tracks.FindActiveByID(ctx, trackID)
	if err != nil {
		return donelog.DoneLogID{}, err
	}
	if track == nil || !track.Active {
		return donelog.DoneLogID{}, fmt.Errorf("track %s not active", trackID.String())
	}

	category, err := h.Categories.FindActiveByID(ctx, categoryID)
	if err != nil {
		return donelog.DoneLogID{}, err
	}
	if category == nil || !category.Active {
		return donelog.DoneLogID{}, fmt.Errorf("category %s not active", categoryID.String())
	}

	id, err := h.IDs.NewDoneLogID(ctx)
	if err != nil {
		return donelog.DoneLogID{}, err
	}

	occurredOn, err := donelog.NewOccurredOn(cmd.OccurredOn)
	if err != nil {
		return donelog.DoneLogID{}, err
	}

	log, err := donelog.NewDoneLog(id, title, trackID, categoryID, count, occurredOn)
	if err != nil {
		return donelog.DoneLogID{}, err
	}

	if err := h.DoneLogs.Save(ctx, log); err != nil {
		return donelog.DoneLogID{}, err
	}

	return id, nil
}
