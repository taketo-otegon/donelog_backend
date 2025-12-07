package donelog

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	ulidPattern = regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}$`)
	slugPattern = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)
)

// DoneLogID represents the identifier of a DONELOG entry.
type DoneLogID struct {
	value string
}

// NewDoneLogID validates and creates a DoneLogID.
func NewDoneLogID(value string) (DoneLogID, error) {
	if value == "" {
		return DoneLogID{}, errors.New("DONELOG id must not be empty")
	}
	if !ulidPattern.MatchString(value) {
		return DoneLogID{}, fmt.Errorf("invalid DONELOG id: %s", value)
	}
	return DoneLogID{value: value}, nil
}

// String returns the string form.
func (id DoneLogID) String() string {
	return id.value
}

// Title represents a DONELOG title.
type Title struct {
	value string
}

const maxTitleLength = 120

// NewTitle validates and creates a Title.
func NewTitle(value string) (Title, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return Title{}, errors.New("title must not be empty")
	}
	if strings.Contains(trimmed, "\n") {
		return Title{}, errors.New("title must not contain line breaks")
	}
	if len([]rune(trimmed)) > maxTitleLength {
		return Title{}, fmt.Errorf("title must be <= %d characters", maxTitleLength)
	}
	return Title{value: trimmed}, nil
}

// String returns the primitive value.
func (t Title) String() string {
	return t.value
}

// TrackID identifies a Track aggregate.
type TrackID struct {
	value string
}

// NewTrackID validates and creates a TrackID.
func NewTrackID(value string) (TrackID, error) {
	if value == "" {
		return TrackID{}, errors.New("track id must not be empty")
	}
	if !slugPattern.MatchString(value) {
		return TrackID{}, fmt.Errorf("invalid track id: %s", value)
	}
	return TrackID{value: value}, nil
}

// String returns the identifier value.
func (id TrackID) String() string {
	return id.value
}

// CategoryID identifies a Category aggregate.
type CategoryID struct {
	value string
}

// NewCategoryID validates and creates a CategoryID.
func NewCategoryID(value string) (CategoryID, error) {
	if value == "" {
		return CategoryID{}, errors.New("category id must not be empty")
	}
	if !slugPattern.MatchString(value) {
		return CategoryID{}, fmt.Errorf("invalid category id: %s", value)
	}
	return CategoryID{value: value}, nil
}

// String returns the identifier value.
func (id CategoryID) String() string {
	return id.value
}

// Count expresses "how many things were done".
type Count struct {
	value int
}

const minCountValue = 1

// NewCount validates and creates a Count greater than zero.
func NewCount(value int) (Count, error) {
	if value < minCountValue {
		return Count{}, fmt.Errorf("count must be >= %d", minCountValue)
	}
	return Count{value: value}, nil
}

// newCountFromNonNegative creates a Count allowing zero (for derived models).
func newCountFromNonNegative(value int) (Count, error) {
	if value < 0 {
		return Count{}, errors.New("count must be >= 0")
	}
	if value == 0 {
		return Count{}, nil
	}
	return Count{value: value}, nil
}

// Int returns the primitive value.
func (c Count) Int() int {
	return c.value
}

// Add returns a new Count after adding delta.
func (c Count) Add(delta Count) (Count, error) {
	result := c.value + delta.value
	return NewCount(result)
}

// Sub returns a new Count after subtracting delta.
func (c Count) Sub(delta Count) (Count, error) {
	result := c.value - delta.value
	return NewCount(result)
}

// OccurredOn represents the date when a DONELOG happened.
type OccurredOn struct {
	date time.Time
}

// NewOccurredOn creates an OccurredOn from YYYY-MM-DD.
func NewOccurredOn(value string) (OccurredOn, error) {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return OccurredOn{}, fmt.Errorf("invalid date: %w", err)
	}
	return OccurredOn{date: t}, nil
}

// OccurredOnFromTime normalizes a time.Time to date precision.
func OccurredOnFromTime(t time.Time) OccurredOn {
	y, m, d := t.Date()
	return OccurredOn{date: time.Date(y, m, d, 0, 0, 0, 0, t.Location())}
}

// Time returns the underlying date as time.Time.
func (o OccurredOn) Time() time.Time {
	return o.date
}

// String returns the YYYY-MM-DD representation.
func (o OccurredOn) String() string {
	return o.date.Format("2006-01-02")
}

// Period represents a closed interval between two dates.
type Period struct {
	start time.Time
	end   time.Time
}

// NewPeriod creates a period (inclusive) from start to end.
func NewPeriod(start, end OccurredOn) (Period, error) {
	if end.Time().Before(start.Time()) {
		return Period{}, errors.New("period end must be on or after start")
	}
	return Period{start: start.Time(), end: end.Time()}, nil
}

// Start returns the start date.
func (p Period) Start() time.Time {
	return p.start
}

// End returns the end date.
func (p Period) End() time.Time {
	return p.end
}

// Contains reports whether the given date falls inside the period.
func (p Period) Contains(o OccurredOn) bool {
	t := o.Time()
	return !t.Before(p.start) && !t.After(p.end)
}
