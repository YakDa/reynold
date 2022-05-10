package reynold

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestAddData(t *testing.T) {
	AddData("key1", PerfData{
		Message:   "This is a test",
		Timestamp: time.Now(),
		Duration:  900,
	})

	assert.Equal(t, 1, Count(), "One logger should be added")
}

func TestOrdering(t *testing.T) {
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

func TestClean(t *testing.T) {
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

	assert.Equal(t, 3, len(perfLoggers.hash), "length should be 3")
	assert.Equal(t, 3, len(perfLoggers.list), "length should be 3")

	Clean()
	assert.Equal(t, 0, len(perfLoggers.hash), "logs should be cleaned")
}
