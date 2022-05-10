package reynold

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestNewPerfSingleton(t *testing.T) {
	AddData("key1", PerfData{
		Message:   "This is a test",
		Timestamp: time.Now(),
		Duration:  900,
	})

	assert.Equal(t, 1, Count(), "One logger should be added")
}

func TestPerfLoggers_Ordering(t *testing.T) {
	AddData("key1", PerfData{
		Message:   "This is a test",
		Timestamp: time.Now(),
		Duration:  900,
	})
	AddData("key2", PerfData{
		Message:   "This is a test",
		Timestamp: time.Now(),
		Duration:  900,
	})
	AddData("key3", PerfData{
		Message:   "This is a test",
		Timestamp: time.Now(),
		Duration:  900,
	})

	assert.Equal(t, "key1", perfLoggers.list[0])
	assert.Equal(t, "key2", perfLoggers.list[1])
	assert.Equal(t, "key3", perfLoggers.list[2])
	log.Println(String())
}
