package sqlite

import (
	"testing"
	"time"
)

func TestTimeFromInterface(t *testing.T) {
	expected := time.Date(
		2026, time.October, 26,
		9, 30, 0, 123456789,
		time.UTC,
	)

	tests := []struct {
		name        string
		value       any
		expected    time.Time
		expectError bool
	}{
		{
			name:     "nil",
			value:    nil,
			expected: time.Time{},
		},
		{
			name:     "time",
			value:    expected,
			expected: expected,
		},
		{
			name:     "time pointer",
			value:    &expected,
			expected: expected,
		},
		{
			name:     "RFC3339 timestamp",
			value:    expected.Format(time.RFC3339Nano),
			expected: expected,
		},
		{
			name:        "invalid timestamp",
			value:       "invalid",
			expectError: true,
		},
		{
			name:        "unsupported type",
			value:       123,
			expectError: true,
		},
		{
			name:     "Go timestamp representation",
			value:    "0001-01-01 00:00:00 +0000 UTC",
			expected: time.Time{},
		},
		{
			name:  "Go timestamp with fractional seconds",
			value: "2026-10-26 09:30:00.123456789 +0000 UTC",
			expected: time.Date(
				2026, time.October, 26,
				9, 30, 0, 123456789,
				time.UTC,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := timeFromInterface(tt.value)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !actual.Equal(tt.expected) {
				t.Errorf(
					"expected %v, got %v",
					tt.expected,
					actual,
				)
			}
		})
	}
}
