package utils

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	Err1 = errors.New("error 1")
	Err2 = errors.New("error 2")
	Err3 = errors.New("error 3")
)

func TestWithRetries(t *testing.T) {
	tests := []struct {
		name                 string
		contextCreationFunc  func() context.Context
		retriesCount         uint
		timeSleep            time.Duration
		functionCreationFunc func() func() error
		wantErr              assert.ErrorAssertionFunc
		expectedErrors       []error
	}{
		{
			name: "cancellation by context",
			contextCreationFunc: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
				return ctx
			},
			retriesCount: 4,
			timeSleep:    time.Millisecond * 50,
			functionCreationFunc: func() func() error {
				return func() error {
					return assert.AnError
				}
			},
			wantErr:        assert.Error,
			expectedErrors: []error{context.Canceled},
		},
		{
			name: "cancellation by context endless",
			contextCreationFunc: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
				return ctx
			},
			retriesCount: 0,
			timeSleep:    time.Millisecond * 10,
			functionCreationFunc: func() func() error {
				return func() error {
					return assert.AnError
				}
			},
			wantErr:        assert.Error,
			expectedErrors: []error{context.Canceled},
		},
		{
			name: "check all errors",
			contextCreationFunc: func() context.Context {
				return context.Background()
			},
			retriesCount: 3,
			timeSleep:    0,
			functionCreationFunc: func() func() error {
				try := 0
				return func() error {
					defer func() {
						try++
					}()
					switch try {
					case 0:
						return Err1
					case 1:
						return Err2
					case 2:
						return Err3
					}
					return errors.New("not implemented error")
				}
			},
			wantErr:        assert.Error,
			expectedErrors: []error{Err1, Err2, Err3},
		},
		{
			name: "endless retries first try",
			contextCreationFunc: func() context.Context {
				return context.Background()
			},
			retriesCount: 0,
			timeSleep:    time.Millisecond * 3,
			functionCreationFunc: func() func() error {
				return func() error {
					return nil
				}
			},
			wantErr:        assert.NoError,
			expectedErrors: nil,
		},
		{
			name: "endless retries big num try",
			contextCreationFunc: func() context.Context {
				return context.Background()
			},
			retriesCount: 0,
			timeSleep:    time.Millisecond * 3,
			functionCreationFunc: func() func() error {
				i := 0
				return func() error {
					if i < 30 {
						i++
						return assert.AnError
					}
					return nil
				}
			},
			wantErr:        assert.NoError,
			expectedErrors: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.contextCreationFunc()
			function := tt.functionCreationFunc()

			got := WithRetries(ctx, tt.retriesCount, tt.timeSleep, function, slog.Default())

			if tt.wantErr(t, got) {
				for _, err := range tt.expectedErrors {
					assert.ErrorIs(t, got, err)
				}
			}
		})
	}
}
