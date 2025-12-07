package command

import "github.com/taketosaeki/donelog/internal/domain/donelog"

// Track is a minimal representation used by commands to validate references.
type Track struct {
	ID              donelog.TrackID
	DefaultCategory *donelog.CategoryID
	Active          bool
}

// Category is a minimal representation used by commands.
type Category struct {
	ID     donelog.CategoryID
	Active bool
}
