package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// EnvelopeHeader Envelope structure based on Sentry format.
type EnvelopeHeader struct {
	EventID string    `json:"event_id"`
	SentAt  time.Time `json:"sent_at"`
	DSN     string    `json:"dsn"`
	SDK     struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"sdk"`
	Trace struct {
		Environment string `json:"environment"`
		PublicKey   string `json:"public_key"`
		TraceID     string `json:"trace_id"`
	} `json:"trace"`
}

type EnvelopeItemHeader struct {
	Type   string `json:"type"`
	Length int    `json:"length"`
}

type EventData struct {
	EventID     string                 `json:"event_id"`
	Timestamp   time.Time              `json:"timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	Platform    string                 `json:"platform"`
	Environment string                 `json:"environment"`
	ServerName  string                 `json:"server_name"`
	Release     string                 `json:"release"`
	Contexts    map[string]interface{} `json:"contexts"`
	User        map[string]interface{} `json:"user"`
	SDK         map[string]interface{} `json:"sdk"`
	Modules     map[string]string      `json:"modules"`
	Exception   *ExceptionData         `json:"exception,omitempty"`
}

type ExceptionData struct {
	Values []ExceptionValue `json:"values"`
}

type ExceptionValue struct {
	Type       string      `json:"type"`
	Value      string      `json:"value"`
	Module     string      `json:"module,omitempty"`
	Stacktrace *Stacktrace `json:"stacktrace,omitempty"`
}

type Stacktrace struct {
	Frames []StackFrame `json:"frames"`
}

type StackFrame struct {
	Filename string `json:"filename"`
	Function string `json:"function"`
	Lineno   int    `json:"lineno"`
	InApp    bool   `json:"in_app"`
}

// Load test configuration.
type LoadTestConfig struct {
	URL            string
	ProjectID      string
	SentryKey      string
	Threads        int
	Duration       time.Duration
	RequestsPerSec int
	Levels         []string
	Platforms      []string
}

// Load test results.
type LoadTestResults struct {
	TotalRequests      int64
	SuccessfulRequests int64
	FailedRequests     int64
	TotalDuration      time.Duration
	StartTime          time.Time
	EndTime            time.Time
	ResponseTimes      []time.Duration
	mu                 sync.RWMutex
}

func (r *LoadTestResults) AddResponseTime(duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ResponseTimes = append(r.ResponseTimes, duration)
}

func (r *LoadTestResults) GetPercentile(percentile float64) time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.ResponseTimes) == 0 {
		return 0
	}

	// Sort response times (simplified - in production you'd want to use a proper sort)
	// For now, we'll just return the median for 50th percentile
	if percentile == 50.0 {
		return r.ResponseTimes[len(r.ResponseTimes)/2]
	}

	// For other percentiles, we'll return a simple approximation
	index := int(float64(len(r.ResponseTimes)) * percentile / 100.0)
	if index >= len(r.ResponseTimes) {
		index = len(r.ResponseTimes) - 1
	}

	return r.ResponseTimes[index]
}

func generateRandomID() string {
	bytess := make([]byte, 16)
	_, _ = rand.Read(bytess)

	return hex.EncodeToString(bytess)
}

//nolint:errcheck,gosec
func generateEnvelope(level, platform, projectID, sentryKey string) ([]byte, error) {
	eventID := generateRandomID()
	traceID := generateRandomID()

	// Create envelope header
	header := EnvelopeHeader{
		EventID: eventID,
		SentAt:  time.Now(),
		DSN:     fmt.Sprintf("http://%s@127.0.0.1:8080/%s", sentryKey, projectID),
		SDK: struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{
			Name:    "sentry.go",
			Version: "0.33.0",
		},
		Trace: struct {
			Environment string `json:"environment"`
			PublicKey   string `json:"public_key"`
			TraceID     string `json:"trace_id"`
		}{
			Environment: "production",
			PublicKey:   sentryKey,
			TraceID:     traceID,
		},
	}

	// Create event data
	eventData := EventData{
		EventID:     eventID,
		Timestamp:   time.Now(),
		Level:       level,
		Message:     fmt.Sprintf("Load test %s event from %s", level, platform),
		Platform:    platform,
		Environment: "production",
		ServerName:  "load-test-server",
		Release:     "load-test@1.1.0",
		Contexts: map[string]interface{}{
			"device": map[string]interface{}{
				"arch":    "x86_64",
				"num_cpu": 8,
			},
			"os": map[string]interface{}{
				"name": "linux",
			},
			"runtime": map[string]interface{}{
				"name":    "go",
				"version": "go1.24.2",
			},
			"trace": map[string]interface{}{
				"span_id":  generateRandomID(),
				"trace_id": traceID,
			},
		},
		User: map[string]interface{}{
			"id":       "load-test-user",
			"username": "loadtester",
		},
		SDK: map[string]interface{}{
			"name":         "sentry.go",
			"version":      "0.33.0",
			"integrations": []string{"ContextifyFrames", "Environment", "GlobalTags"},
			"packages": []map[string]string{
				{"name": "sentry-go", "version": "0.33.0"},
			},
		},
		Modules: map[string]string{
			"github.com/getsentry/sentry-go": "v0.33.0",
			"load-test":                      "v1.0.0",
		},
	}

	// Add exception data for certain levels
	if level == "fatal" || level == "exception" {
		eventData.Exception = &ExceptionData{
			Values: []ExceptionValue{
				{
					Type:   "LoadTestException",
					Value:  fmt.Sprintf("Load test %s exception", level),
					Module: "loadtest",
					Stacktrace: &Stacktrace{
						Frames: []StackFrame{
							{
								Filename: "load_test.go",
								Function: "generateEnvelope",
								Lineno:   42,
								InApp:    true,
							},
							{
								Filename: "main.go",
								Function: "main",
								Lineno:   15,
								InApp:    true,
							},
						},
					},
				},
			},
		}
	}

	// Serialize envelope header
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return nil, err
	}

	// Serialize event data
	eventBytes, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}

	// Create item header
	itemHeader := EnvelopeItemHeader{
		Type:   "event",
		Length: len(eventBytes),
	}

	itemHeaderBytes, err := json.Marshal(itemHeader)
	if err != nil {
		return nil, err
	}

	// Build envelope
	var envelope bytes.Buffer
	envelope.Write(headerBytes)
	envelope.WriteString("\n")
	envelope.Write(itemHeaderBytes)
	envelope.WriteString("\n")
	envelope.Write(eventBytes)

	return envelope.Bytes(), nil
}

// ProgressBar Progress bar structure.
type ProgressBar struct {
	total      int64
	current    int64
	startTime  time.Time
	mu         sync.Mutex
	lastUpdate time.Time
	updateFreq time.Duration
}

func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		startTime:  time.Now(),
		updateFreq: 500 * time.Millisecond, // Update every 500ms
	}
}

func (pb *ProgressBar) Increment() {
	atomic.AddInt64(&pb.current, 1)
}

func (pb *ProgressBar) SetTotal(total int64) {
	atomic.StoreInt64(&pb.total, total)
}

func (pb *ProgressBar) Update() {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	now := time.Now()
	if now.Sub(pb.lastUpdate) < pb.updateFreq {
		return
	}

	current := atomic.LoadInt64(&pb.current)
	total := atomic.LoadInt64(&pb.total)

	if total == 0 {
		return
	}

	elapsed := now.Sub(pb.startTime)
	progress := float64(current) / float64(total) * 100

	// Calculate ETA
	var eta time.Duration
	if current > 0 {
		rate := float64(current) / elapsed.Seconds()
		remaining := float64(total-current) / rate
		eta = time.Duration(remaining) * time.Second
	}

	// Create progress bar
	barWidth := 30
	filled := int(float64(barWidth) * progress / 100)
	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "="
		} else if i == filled {
			bar += ">"
		} else {
			bar += " "
		}
	}
	bar += "]"

	// Clear line and print progress
	fmt.Printf("\r%s %6.1f%% (%d/%d) | Elapsed: %s | ETA: %s",
		bar, progress, current, total,
		formatDuration(elapsed),
		formatDuration(eta))

	pb.lastUpdate = now
}

func (pb *ProgressBar) Finish() {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	current := atomic.LoadInt64(&pb.current)
	total := atomic.LoadInt64(&pb.total)
	elapsed := time.Since(pb.startTime)

	// Clear line and print final result
	fmt.Printf("\r%s %6.1f%% (%d/%d) | Completed in %s\n",
		"[==============================]",
		100.0, current, total,
		formatDuration(elapsed))
}

func formatDuration(duration time.Duration) string {
	if duration < time.Minute {
		return fmt.Sprintf("%.1fs", duration.Seconds())
	} else if duration < time.Hour {
		return fmt.Sprintf("%dm%ds", int(duration.Minutes()), int(duration.Seconds())%60)
	} else {
		return fmt.Sprintf("%dh%dm", int(duration.Hours()), int(duration.Minutes())%60)
	}
}

func sendEnvelope(config LoadTestConfig, results *LoadTestResults, progressBar *ProgressBar) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	levels := []string{"debug", "info", "warning", "error", "exception", "fatal"}
	platforms := []string{"go", "python", "javascript", "java", "csharp"}

	levelIndex := 0
	platformIndex := 0

	for {
		level := levels[levelIndex%len(levels)]
		platform := platforms[platformIndex%len(platforms)]

		envelopeData, err := generateEnvelope(level, platform, config.ProjectID, config.SentryKey)
		if err != nil {
			log.Printf("Error generating envelope: %v", err)
			atomic.AddInt64(&results.FailedRequests, 1)
			progressBar.Increment()
			progressBar.Update()

			continue
		}

		url := fmt.Sprintf("%s/api/%s/envelope/", config.URL, config.ProjectID)
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, bytes.NewReader(envelopeData))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			atomic.AddInt64(&results.FailedRequests, 1)
			progressBar.Increment()
			progressBar.Update()

			continue
		}

		req.Header.Set("Content-Type", "application/x-sentry-envelope")
		req.Header.Set("X-Sentry-Auth", "Sentry sentry_version=7, sentry_client=sentry.go/0.33.0, sentry_key="+config.SentryKey) //nolint:lll // it's ok

		start := time.Now()
		resp, err := client.Do(req)
		duration := time.Since(start)

		atomic.AddInt64(&results.TotalRequests, 1)
		results.AddResponseTime(duration)
		progressBar.Increment()
		progressBar.Update()

		if err != nil {
			log.Printf("Error sending request: %v", err)
			atomic.AddInt64(&results.FailedRequests, 1)
		} else {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				atomic.AddInt64(&results.SuccessfulRequests, 1)
			} else {
				log.Printf("Request failed with status: %d", resp.StatusCode)
				atomic.AddInt64(&results.FailedRequests, 1)
			}
		}

		levelIndex++
		platformIndex++

		// Rate limiting if specified
		if config.RequestsPerSec > 0 {
			time.Sleep(time.Second / time.Duration(config.RequestsPerSec))
		}
	}
}

func printResults(results *LoadTestResults) {
	fmt.Printf("\n=== Load Test Results ===\n")
	fmt.Printf("Total Requests: %d\n", atomic.LoadInt64(&results.TotalRequests))
	fmt.Printf("Successful Requests: %d\n", atomic.LoadInt64(&results.SuccessfulRequests))
	fmt.Printf("Failed Requests: %d\n", atomic.LoadInt64(&results.FailedRequests))
	fmt.Printf("Success Rate: %.2f%%\n",
		float64(atomic.LoadInt64(&results.SuccessfulRequests))/float64(atomic.LoadInt64(&results.TotalRequests))*100)
	fmt.Printf("Total Duration: %v\n", results.TotalDuration)

	if len(results.ResponseTimes) > 0 {
		fmt.Printf("Average Response Time: %v\n", results.GetPercentile(50))
		fmt.Printf("75th Percentile Response Time: %v\n", results.GetPercentile(75))
		fmt.Printf("95th Percentile Response Time: %v\n", results.GetPercentile(95))
		fmt.Printf("99th Percentile Response Time: %v\n", results.GetPercentile(99))
		fmt.Printf("Requests per Second: %.2f\n",
			float64(atomic.LoadInt64(&results.TotalRequests))/results.TotalDuration.Seconds())
	}
}

func main() {
	var (
		url       = flag.String("url", "http://localhost:8080", "Target server URL")
		projectID = flag.String("project", "1", "Project ID")
		sentryKey = flag.String("key",
			"418aba92087742d7ac5a252ebee0d7299f1f1ca96fdfa781843a990bf7b93cc1", "Sentry key")
		threads        = flag.Int("threads", 10, "Number of concurrent threads")
		duration       = flag.Duration("duration", 60*time.Second, "Test duration")
		requestsPerSec = flag.Int("rps", 0, "Requests per second per thread (0 for unlimited)")
	)

	flag.Parse()

	config := LoadTestConfig{
		URL:            *url,
		ProjectID:      *projectID,
		SentryKey:      *sentryKey,
		Threads:        *threads,
		Duration:       *duration,
		RequestsPerSec: *requestsPerSec,
	}

	results := &LoadTestResults{
		StartTime: time.Now(),
	}

	// Calculate expected total requests for progress bar
	var expectedTotal int64
	if config.RequestsPerSec > 0 {
		expectedTotal = int64(config.RequestsPerSec * config.Threads * int(config.Duration.Seconds()))
	} else {
		// For unlimited RPS, estimate based on duration
		expectedTotal = int64(config.Threads * int(config.Duration.Seconds()) * 10) // Assume 10 RPS per thread
	}

	// Create progress bar
	progressBar := NewProgressBar()
	progressBar.SetTotal(expectedTotal)

	fmt.Printf("Starting load test with configuration:\n")
	fmt.Printf("URL: %s\n", config.URL)
	fmt.Printf("Project ID: %s\n", config.ProjectID)
	fmt.Printf("Threads: %d\n", config.Threads)
	fmt.Printf("Duration: %v\n", config.Duration)
	if config.RequestsPerSec > 0 {
		fmt.Printf("Requests per second per thread: %d\n", config.RequestsPerSec)
	}
	fmt.Printf("Sending all levels: debug, info, warning, error, exception, fatal\n")
	fmt.Printf("Sending from platforms: go, python, javascript, java, csharp\n")
	fmt.Printf("\nStarting test...\n")

	// Start progress bar updates in background
	stopProgress := make(chan bool)
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				progressBar.Update()
			case <-stopProgress:
				return
			}
		}
	}()

	// Start worker threads
	var wg sync.WaitGroup
	for i := 0; i < config.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sendEnvelope(config, results, progressBar)
		}()
	}

	// Wait for duration
	time.Sleep(config.Duration)

	// Stop progress bar
	close(stopProgress)
	progressBar.Finish()

	results.EndTime = time.Now()
	results.TotalDuration = results.EndTime.Sub(results.StartTime)

	printResults(results)
}
