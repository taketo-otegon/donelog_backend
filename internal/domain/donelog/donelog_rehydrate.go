package donelog

import "time"

// RawDoneLog represents persisted primitive values for rehydration.
type RawDoneLog struct {
	ID         string
	Title      string
	TrackID    string
	CategoryID string
	Count      int
	OccurredOn time.Time
}

// RehydrateDoneLog rebuilds a DoneLog aggregate from persisted primitives.
func RehydrateDoneLog(raw RawDoneLog) (*DoneLog, error) {
	id, err := NewDoneLogID(raw.ID)
	if err != nil {
		return nil, err
	}
	title, err := NewTitle(raw.Title)
	if err != nil {
		return nil, err
	}
	trackID, err := NewTrackID(raw.TrackID)
	if err != nil {
		return nil, err
	}
	categoryID, err := NewCategoryID(raw.CategoryID)
	if err != nil {
		return nil, err
	}
	count, err := NewCount(raw.Count)
	if err != nil {
		return nil, err
	}
	occurredOn := OccurredOnFromTime(raw.OccurredOn)

	return NewDoneLog(id, title, trackID, categoryID, count, occurredOn)
}
