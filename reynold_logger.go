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

type PerfLoggers map[string][]PerfData

var (
	once        sync.Once
	perfLoggers PerfLoggers
	mu          sync.Mutex
)

func NewPerfSingleton() *PerfLoggers {
	once.Do(func() {
		perfLoggers = make(map[string][]PerfData)
	})
	return &perfLoggers
}

func (p *PerfLoggers) AddPerfData(key, message string) func() {
	start := time.Now()
	ts := PerfData{
		Message:   message,
		Timestamp: start,
	}
	return func() {
		ts.Duration = time.Since(start)
		p.AddData(key, ts)
	}
}

func (p *PerfLoggers) AddData(key string, data PerfData) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := (*p)[key]; ok {
		(*p)[key] = append((*p)[key], data)
		return
	}
	(*p)[key] = []PerfData{data}
}

func (p *PerfLoggers) String() string {
	mu.Lock()
	defer mu.Unlock()
	result := "\n[Benchmark Start] \n"
	for s, logger := range *p {
		result += fmt.Sprintf("\n------------------Logger %s----------------------\n", s)
		var execTime time.Duration = 0
		prevIndex := 0
		for j, perf := range logger {
			if !strings.HasPrefix(s, "[inner]") {
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
		if strings.HasPrefix(s, "[inner]") {
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

func (p *PerfLoggers) Clean() {
	mu.Lock()
	defer mu.Unlock()
	perfLoggers = make(PerfLoggers)
}
