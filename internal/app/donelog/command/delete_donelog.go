package command

import (
	"context"
	"fmt"

	"github.com/taketosaeki/donelog/internal/domain/donelog"
)

// DeleteDoneLogCommand removes a DONELOG entry.
type DeleteDoneLogCommand struct {
	ID string
}

func (c DeleteDoneLogCommand) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// DeleteDoneLogHandler handles DeleteDoneLogCommand.
type DeleteDoneLogHandler struct {
	DoneLogs DoneLogRepository
}

func (h DeleteDoneLogHandler) Handle(ctx context.Context, cmd DeleteDoneLogCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	id, err := donelog.NewDoneLogID(cmd.ID)
	if err != nil {
		return err
	}

	return h.DoneLogs.Delete(ctx, id)
}
