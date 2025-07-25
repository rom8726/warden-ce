package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvent_FullFingerprint(t *testing.T) {
	ev := Event{
		Message:  "some message",
		Level:    IssueLevelError,
		Platform: "go",
	}

	exType := "some type"
	ev.ExceptionType = &exType

	exValue := "some value"
	ev.ExceptionValue = &exValue

	exStackTrace := []byte("some stacktrace")
	ev.ExceptionStacktrace = exStackTrace

	t.Run("event source", func(t *testing.T) {
		ev.Source = SourceEvent
		assert.Equal(t, "d007bc4f570309a35f10c47ddf30edae9af853a8", ev.FullFingerprint())
	})
	t.Run("exception source", func(t *testing.T) {
		ev.Source = SourceException
		assert.Equal(t, "a3da4ed0ab522704072d2ffef63601ad54fbde89", ev.FullFingerprint())
	})
}
