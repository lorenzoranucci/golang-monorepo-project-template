package retry

import (
	"errors"
	"testing"
	"time"
)

func Test_withExponentialBackoff(t *testing.T) {
	tests := []struct {
		name           string
		funcToRetry    func() error
		initialDelay   time.Duration
		maxWait        time.Duration
		mockSleep      *mockSleep
		wantSleepCalls int
		wantErr        bool
	}{
		{
			name: "WHEN funcToRetry returns nil THEN withExponentialBackoff returns nil",
			funcToRetry: func() error {
				return nil
			},
			mockSleep:    &mockSleep{},
			maxWait:      300 * time.Millisecond,
			initialDelay: 100 * time.Millisecond,
			wantErr:      false,
		},
		{
			name: "WHEN funcToRetry returns error " +
				"THEN withExponentialBackoff returns always error after retrying maxWait time",
			funcToRetry: func() error {
				return errors.New("error")
			},
			mockSleep:    &mockSleep{},
			maxWait:      300 * time.Millisecond,
			initialDelay: 100 * time.Millisecond,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := withExponentialBackoff(
				tt.funcToRetry,
				tt.initialDelay,
				tt.maxWait,
				tt.mockSleep.Sleep,
			); (err != nil) != tt.wantErr {
				t.Errorf("withExponentialBackoff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type mockSleep struct{}

func (m *mockSleep) Sleep(delay time.Duration) {
}
