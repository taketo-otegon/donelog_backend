package donelog

import (
	"testing"
	"time"
)

func TestNewDoneLogID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"OK: valid ULID", "01HYR1X5C9XM9P6H7K71M9QAHX", false},
		{"NG: empty", "", true},
		{"NG: bad format", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewDoneLogID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && id.String() != tt.input {
				t.Fatalf("unexpected id: %s", id.String())
			}
		})
	}
}

func TestNewTitle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantVal string
	}{
		{"OK: trims surrounding spaces", "  Study Go  ", false, "Study Go"},
		{"NG: empty title", "", true, ""},
		{"NG: contains newline", "Hello\nWorld", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, err := NewTitle(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && title.String() != tt.wantVal {
				t.Fatalf("expected %q, got %q", tt.wantVal, title.String())
			}
		})
	}
}

func TestNewTrackID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"OK: valid slug", "track_sample", false},
		{"NG: empty", "", true},
		{"NG: uppercase not allowed", "Track", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewTrackID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && id.String() != tt.input {
				t.Fatalf("unexpected track id: %s", id.String())
			}
		})
	}
}

func TestNewCategoryID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"OK: valid slug", "cat_sample", false},
		{"NG: empty", "", true},
		{"NG: uppercase not allowed", "Cat", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewCategoryID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && id.String() != tt.input {
				t.Fatalf("unexpected category id: %s", id.String())
			}
		})
	}
}

func TestNewCount(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{"OK: positive", 5, false},
		{"NG: zero", 0, true},
		{"NG: negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := NewCount(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && count.Int() != tt.input {
				t.Fatalf("expected %d, got %d", tt.input, count.Int())
			}
		})
	}
}

func TestOccurredOn(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantStr string
	}{
		{"OK: valid date", "2024-05-01", false, "2024-05-01"},
		{"NG: invalid format", "2024/05/01", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := NewOccurredOn(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil && o.String() != tt.wantStr {
				t.Fatalf("expected %s, got %s", tt.wantStr, o.String())
			}
		})
	}

	t.Run("normalize from time", func(t *testing.T) {
		fromTime := OccurredOnFromTime(time.Date(2023, 12, 10, 10, 0, 0, 0, time.UTC))
		if fromTime.String() != "2023-12-10" {
			t.Fatalf("expected normalized date, got %s", fromTime.String())
		}
	})
}

func TestPeriodContains(t *testing.T) {
	start, _ := NewOccurredOn("2024-01-01")
	end, _ := NewOccurredOn("2024-01-10")
	period, err := NewPeriod(start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name   string
		input  string
		wantIn bool
	}{
		{"inside", "2024-01-05", true},
		{"start boundary", "2024-01-01", true},
		{"end boundary", "2024-01-10", true},
		{"outside", "2023-12-31", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, _ := NewOccurredOn(tt.input)
			if got := period.Contains(date); got != tt.wantIn {
				t.Fatalf("Contains(%s) = %v, want %v", tt.input, got, tt.wantIn)
			}
		})
	}
}

func TestNewPeriod(t *testing.T) {
	tests := []struct {
		name    string
		start   string
		end     string
		wantErr bool
	}{
		{"OK: start before end", "2024-01-01", "2024-01-10", false},
		{"NG: start after end", "2024-01-10", "2024-01-01", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, _ := NewOccurredOn(tt.start)
			end, _ := NewOccurredOn(tt.end)
			_, err := NewPeriod(start, end)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
