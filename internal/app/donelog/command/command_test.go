package command

import (
	"context"
	"testing"
	"time"

	"github.com/taketosaeki/donelog/internal/domain/donelog"
)

type mockDoneLogRepo struct {
	saved *donelog.DoneLog
	found *donelog.RawDoneLog
	err   error
}

func (m *mockDoneLogRepo) Save(ctx context.Context, log *donelog.DoneLog) error {
	m.saved = log
	return m.err
}
func (m *mockDoneLogRepo) FindByID(ctx context.Context, id donelog.DoneLogID) (*donelog.RawDoneLog, error) {
	return m.found, m.err
}
func (m *mockDoneLogRepo) Delete(ctx context.Context, id donelog.DoneLogID) error {
	return m.err
}

type mockTrackRepo struct {
	track *Track
	err   error
}

func (m mockTrackRepo) FindActiveByID(ctx context.Context, id donelog.TrackID) (*Track, error) {
	return m.track, m.err
}

type mockCategoryRepo struct {
	category *Category
	err      error
}

func (m mockCategoryRepo) FindActiveByID(ctx context.Context, id donelog.CategoryID) (*Category, error) {
	return m.category, m.err
}

type mockIDGenerator struct {
	id  donelog.DoneLogID
	err error
}

func (m mockIDGenerator) NewDoneLogID(ctx context.Context) (donelog.DoneLogID, error) {
	return m.id, m.err
}

func TestCreateDoneLog(t *testing.T) {
	tests := []struct {
		name     string
		cmd      CreateDoneLogCommand
		track    *Track
		category *Category
		idGenErr error
		wantErr  bool
		wantID   string
	}{
		{
			name: "OK: create done log",
			cmd: CreateDoneLogCommand{
				Title:      "Test",
				TrackID:    "track_sample",
				CategoryID: "cat_sample",
				Count:      2,
				OccurredOn: "2024-05-01",
			},
			track:    &Track{Active: true},
			category: &Category{Active: true},
			wantID:   "01HYR1X5C9XM9P6H7K71M9QAHX",
		},
		{
			name: "NG: inactive track",
			cmd: CreateDoneLogCommand{
				Title:      "Test",
				TrackID:    "track_sample",
				CategoryID: "cat_sample",
				Count:      2,
				OccurredOn: "2024-05-01",
			},
			track:    &Track{Active: false},
			category: &Category{Active: true},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockDoneLogRepo{}
			handler := CreateDoneLogHandler{
				DoneLogs:   repo,
				Tracks:     mockTrackRepo{track: tt.track},
				Categories: mockCategoryRepo{category: tt.category},
				IDs:        mockIDGenerator{id: mustDoneLogID(t, "01HYR1X5C9XM9P6H7K71M9QAHX"), err: tt.idGenErr},
			}

			id, err := handler.Handle(context.Background(), tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && id.String() != tt.wantID {
				t.Fatalf("unexpected id: %s", id.String())
			}
		})
	}
}

func TestUpdateDoneLog_NotFound(t *testing.T) {
	tests := []struct {
		name     string
		found    *donelog.RawDoneLog
		category *Category
		wantErr  bool
	}{
		{
			name:     "NG: missing log",
			found:    nil,
			category: &Category{Active: true},
			wantErr:  true,
		},
		{
			name: "NG: inactive category",
			found: &donelog.RawDoneLog{
				ID:         "01HYR1X5C9XM9P6H7K71M9QAHX",
				Title:      "Existing",
				TrackID:    "track_sample",
				CategoryID: "cat_old",
				Count:      1,
				OccurredOn: donelog.OccurredOnFromTime(time.Now()).Time(),
			},
			category: &Category{Active: false},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockDoneLogRepo{found: tt.found}
			handler := UpdateDoneLogHandler{
				DoneLogs: repo,
				Categories: mockCategoryRepo{
					category: tt.category,
				},
			}

			cmd := UpdateDoneLogCommand{
				ID:         "01HYR1X5C9XM9P6H7K71M9QAHX",
				Title:      "Updated",
				CategoryID: "cat_sample",
				Count:      5,
				OccurredOn: "2024-05-02",
			}

			err := handler.Handle(context.Background(), cmd)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteDoneLog(t *testing.T) {
	tests := []struct {
		name    string
		repoErr error
		wantErr bool
	}{
		{"OK: delete success", nil, false},
		{"NG: repo error", assertErr("delete failed"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockDoneLogRepo{err: tt.repoErr}
			handler := DeleteDoneLogHandler{DoneLogs: repo}

			cmd := DeleteDoneLogCommand{ID: "01HYR1X5C9XM9P6H7K71M9QAHX"}

			err := handler.Handle(context.Background(), cmd)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

// helper
type mockClock struct {
	value donelog.OccurredOn
	err   error
}

type assertErr string

func (e assertErr) Error() string { return string(e) }

func mustDoneLogID(t *testing.T, value string) donelog.DoneLogID {
	t.Helper()
	id, err := donelog.NewDoneLogID(value)
	if err != nil {
		t.Fatalf("failed to create DoneLogID: %v", err)
	}
	return id
}
