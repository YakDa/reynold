package reynold

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestNewPerfSingleton(t *testing.T) {
	NewPerfSingleton().AddData("key1", PerfData{
		Message:   "This is a test",
		Timestamp: time.Now(),
		Duration:  900,
	})

	assert.Equal(t, 1, NewPerfSingleton().Count(), "One logger should be added")
}

func TestPerfLoggers_Ordering(t *testing.T) {
	NewPerfSingleton().AddData("key1", PerfData{
		Message:   "This is a test",
		Timestamp: time.Now(),
		Duration:  900,
	})
	NewPerfSingleton().AddData("key2", PerfData{
		Message:   "This is a test",
		Timestamp: time.Now(),
		Duration:  900,
	})
	NewPerfSingleton().AddData("key3", PerfData{
		Message:   "This is a test",
		Timestamp: time.Now(),
		Duration:  900,
	})

	assert.Equal(t, "key1", NewPerfSingleton().list[0])
	assert.Equal(t, "key2", NewPerfSingleton().list[1])
	assert.Equal(t, "key3", NewPerfSingleton().list[2])
	log.Println(NewPerfSingleton().String())
}
