package reynold

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewPerfSingleton(t *testing.T) {
	NewPerfSingleton().AddData("key1", PerfData{
		Message:   "This is a test",
		Timestamp: time.Now(),
		Duration:  900,
	})

	assert.Equal(t, 1, len(*NewPerfSingleton()), "One logger should be added")
}
