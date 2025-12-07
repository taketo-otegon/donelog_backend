package command

import (
	"context"
	"fmt"

	"github.com/taketosaeki/donelog/internal/domain/donelog"
)

// UpdateDoneLogCommand updates an existing DONELOG.
type UpdateDoneLogCommand struct {
	ID         string
	Title      string
	CategoryID string
	Count      int
	OccurredOn string
}

func (c UpdateDoneLogCommand) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("id is required")
	}
	if c.Title == "" {
		return fmt.Errorf("title is required")
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

// UpdateDoneLogHandler handles UpdateDoneLogCommand.
type UpdateDoneLogHandler struct {
	DoneLogs   DoneLogRepository
	Categories CategoryRepository
}

func (h UpdateDoneLogHandler) Handle(ctx context.Context, cmd UpdateDoneLogCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	id, err := donelog.NewDoneLogID(cmd.ID)
	if err != nil {
		return err
	}
	categoryID, err := donelog.NewCategoryID(cmd.CategoryID)
	if err != nil {
		return err
	}

	rawLog, err := h.DoneLogs.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if rawLog == nil {
		return fmt.Errorf("doneLog %s not found", id.String())
	}

	log, err := donelog.RehydrateDoneLog(*rawLog)
	if err != nil {
		return err
	}

	category, err := h.Categories.FindActiveByID(ctx, categoryID)
	if err != nil {
		return err
	}
	if category == nil || !category.Active {
		return fmt.Errorf("category %s not active", categoryID.String())
	}

	occurredOn, err := donelog.NewOccurredOn(cmd.OccurredOn)
	if err != nil {
		return err
	}

	raw := donelog.RawDoneLog{
		ID:         rawLog.ID,
		Title:      cmd.Title,
		TrackID:    rawLog.TrackID,
		CategoryID: cmd.CategoryID,
		Count:      cmd.Count,
		OccurredOn: occurredOn.Time(),
	}

	log, err = donelog.RehydrateDoneLog(raw)
	if err != nil {
		return err
	}

	return h.DoneLogs.Save(ctx, log)
}
