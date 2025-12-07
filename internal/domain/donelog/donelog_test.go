package donelog

import (
	"testing"
	"time"
)

func TestNewDoneLog(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		title     string
		track     string
		category  string
		count     int
		occurred  string
		wantTrack string
		wantErr   bool
	}{
		{
			name:      "OK: valid aggregate",
			id:        "01HYR1X5C9XM9P6H7K71M9QAHX",
			title:     "Initial",
			track:     "track_sample",
			category:  "cat_default",
			count:     2,
			occurred:  "2024-05-01",
			wantTrack: "track_sample",
			wantErr:   false,
		},
		{
			name:     "NG: invalid title",
			id:       "01HYR1X5C9XM9P6H7K71M9QAHX",
			title:    "",
			track:    "track_sample",
			category: "cat_default",
			count:    2,
			occurred: "2024-05-01",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewDoneLogID(tt.id)
			if err != nil {
				t.Fatalf("id creation failed: %v", err)
			}
			title, err := NewTitle(tt.title)
			if (err != nil) != tt.wantErr {
				t.Fatalf("title error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			trackID, _ := NewTrackID(tt.track)
			categoryID, _ := NewCategoryID(tt.category)
			count, _ := NewCount(tt.count)
			occurredOn, _ := NewOccurredOn(tt.occurred)

			log, err := NewDoneLog(id, title, trackID, categoryID, count, occurredOn)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewDoneLog error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if log.TrackID().String() != tt.wantTrack {
				t.Fatalf("expected trackID %s, got %s", tt.wantTrack, log.TrackID().String())
			}
		})
	}
}

func TestDoneLogUpdate(t *testing.T) {
	tests := []struct {
		name      string
		newTitle  string
		newCat    string
		newCount  int
		newDate   string
		wantTitle string
		wantCat   string
		wantCount int
		wantDate  string
	}{
		{
			name:      "updates mutable fields",
			newTitle:  "Updated Title",
			newCat:    "cat_new",
			newCount:  5,
			newDate:   "2024-05-03",
			wantTitle: "Updated Title",
			wantCat:   "cat_new",
			wantCount: 5,
			wantDate:  "2024-05-03",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := mustDoneLogID(t, "01HYR1X5C9XM9P6H7K71M9QAHX")
			title, _ := NewTitle("Initial")
			trackID, _ := NewTrackID("track_sample")
			categoryID, _ := NewCategoryID("cat_default")
			count, _ := NewCount(2)
			occurredOn, _ := NewOccurredOn("2024-05-01")

			log, err := NewDoneLog(id, title, trackID, categoryID, count, occurredOn)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			newTitle, _ := NewTitle(tt.newTitle)
			newCategory, _ := NewCategoryID(tt.newCat)
			newCount, _ := NewCount(tt.newCount)
			newDate, _ := NewOccurredOn(tt.newDate)

			log.Update(newTitle, newCategory, newCount, newDate)

			if log.Title().String() != tt.wantTitle {
				t.Fatalf("expected title %s, got %s", tt.wantTitle, log.Title())
			}
			if log.CategoryID().String() != tt.wantCat {
				t.Fatalf("expected category %s, got %s", tt.wantCat, log.CategoryID())
			}
			if log.Count().Int() != tt.wantCount {
				t.Fatalf("expected count %d, got %d", tt.wantCount, log.Count().Int())
			}
			if log.OccurredOn().String() != tt.wantDate {
				t.Fatalf("expected date %s, got %s", tt.wantDate, log.OccurredOn())
			}
		})
	}
}

func TestRehydrateDoneLog(t *testing.T) {
	tests := []struct {
		name    string
		raw     RawDoneLog
		wantErr bool
	}{
		{
			name: "OK: restores aggregate",
			raw: RawDoneLog{
				ID:         "01HYR1X5C9XM9P6H7K71M9QAHX",
				Title:      "Rehydrated",
				TrackID:    "track_sample",
				CategoryID: "cat_default",
				Count:      3,
				OccurredOn: time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "NG: invalid raw data",
			raw: RawDoneLog{
				ID:         "invalid-ulid",
				Title:      "",
				TrackID:    "track_sample",
				CategoryID: "cat_default",
				Count:      0,
				OccurredOn: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, err := RehydrateDoneLog(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if log.ID().String() != tt.raw.ID {
				t.Fatalf("expected id %s, got %s", tt.raw.ID, log.ID().String())
			}
			if log.OccurredOn().String() != "2024-05-01" {
				t.Fatalf("expected occurredOn %s, got %s", "2024-05-01", log.OccurredOn().String())
			}
		})
	}
}

func mustDoneLogID(t *testing.T, value string) DoneLogID {
	t.Helper()
	id, err := NewDoneLogID(value)
	if err != nil {
		t.Fatalf("failed to create DoneLogID: %v", err)
	}
	return id
}
