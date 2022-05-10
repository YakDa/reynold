package reynold

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type PerfData struct {
	Message   string
	Timestamp time.Time
	Duration  time.Duration
}

type PerfLoggers struct {
	hash map[string][]PerfData
	list []string
}

var (
	perfLoggers PerfLoggers
	mu          sync.Mutex
)

func init() {
	perfLoggers = PerfLoggers{
		hash: make(map[string][]PerfData),
		list: []string{},
	}
}

func AddPerfData(key, message string) func() {
	start := time.Now()
	ts := PerfData{
		Message:   message,
		Timestamp: start,
	}
	return func() {
		ts.Duration = time.Since(start)
		AddData(key, ts)
	}
}

func AddData(key string, data PerfData) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := perfLoggers.hash[key]; ok {
		perfLoggers.hash[key] = append(perfLoggers.hash[key], data)
		return
	}
	perfLoggers.list = append(perfLoggers.list, key)
	perfLoggers.hash[key] = []PerfData{data}
}

func String() string {
	mu.Lock()
	defer mu.Unlock()
	result := "\n[Benchmark Start] \n"
	for _, key := range perfLoggers.list {
		logger:= perfLoggers.hash[key]
		result += fmt.Sprintf("\n------------------Logger %s----------------------\n", key)
		var execTime time.Duration = 0
		prevIndex := 0
		for j, perf := range logger {
			if !strings.HasPrefix(key, "[inner]") {
				execTime += perf.Duration
				if j != 0 {
					prevPerf := logger[prevIndex]
					prevEnd := prevPerf.Timestamp.Add(prevPerf.Duration)
					duration := perf.Timestamp.Sub(prevEnd)
					if duration < 0 {
						duration = 0
					}
					prevIndex = j
					result += "gap---" + fmt.Sprintf("\t Duration: %v\n", duration)
				}
			}
			result += perf.Message + "---" + "\t  Timestamp:" + perf.Timestamp.Format(time.StampMilli) + "  " + fmt.Sprintf("Duration: %v", perf.Duration) + "\n"
		}
		if strings.HasPrefix(key, "[inner]") {
			continue
		}
		result += "\nSummary:\n"
		var total time.Duration = 0
		if len(logger) > 0 {
			start := logger[0].Timestamp
			lastPerf := logger[len(logger)-1]
			end := lastPerf.Timestamp
			total = end.Sub(start) + lastPerf.Duration
		}
		result += "Total Time: " + total.String() + "\n"
		result += "Internal Execution Time: " + execTime.String() + "\n"
		result += "Network and External Time: " + (total - execTime).String() + "\n"

	}
	result += "[Benchmark End] \n"
	return result
}

func Clean() {
	mu.Lock()
	defer mu.Unlock()
	perfLoggers.hash = make(map[string][]PerfData, 0)
	perfLoggers.list = []string{}
}

func Count() int {
	if len(perfLoggers.hash) != len(perfLoggers.list) {
		return -1
	}
	return len(perfLoggers.list)
}