package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Prometheus metrics for the agent.
	injectionCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "agent_injection_total",
			Help: "Total number of injection requests handled by the agent",
		},
	)
	injectionDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "agent_injection_duration_seconds",
			Help:    "Histogram of fault injection durations",
			Buckets: prometheus.LinearBuckets(5, 5, 10),
		},
	)
)

func init() {
	prometheus.MustRegister(injectionCounter)
	prometheus.MustRegister(injectionDuration)
}

// InjectionRequest represents a fault injection command.
type InjectionRequest struct {
	ExperimentType string `json:"experiment_type"`
	Duration       int    `json:"duration"`     // Duration in seconds
	DelayMs        int    `json:"delay_ms"`     // For network latency
	LossPercent    int    `json:"loss_percent"` // For network packet loss
	CPUWorkers     int    `json:"cpu_workers"`  // For CPU stress
	MemSizeMB      int    `json:"mem_size_mb"`  // For Memory stress
	KillProcess    string `json:"kill_process"` // Pattern or process name to kill (optional)
}

// registerAgentHandlers sets up the agent's HTTP endpoints.
func registerAgentHandlers() {
	http.HandleFunc("/inject", injectHandler)
}

// injectHandler listens for fault injection commands.
func injectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request", http.StatusBadRequest)
		return
	}
	var injReq InjectionRequest
	if err := json.Unmarshal(body, &injReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("[Agent] Received injection request: %+v", injReq)

	// Run the simulation in a separate goroutine so we don't block the HTTP handler.
	go simulateInjection(injReq)

	response := map[string]string{
		"status":  "injected",
		"message": "Fault injection in progress",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// simulateInjection simulates a fault injection experiment.
func simulateInjection(req InjectionRequest) {
	log.Printf("[Agent] Starting experiment: %s for %d seconds...", req.ExperimentType, req.Duration)

	switch req.ExperimentType {
	case "network-latency":
		addNetworkLatency(req.DelayMs)
		time.Sleep(time.Duration(req.Duration) * time.Second)
		removeNetworkLatency()
	case "network-loss":
		addNetworkLoss(req.LossPercent)
		time.Sleep(time.Duration(req.Duration) * time.Second)
		removeNetworkLatency() // same cleanup for netem
	case "cpu-stress":
		runStressNGCPU(req.CPUWorkers, req.Duration)
	case "mem-stress":
		runStressNGMemory(req.MemSizeMB, req.Duration)
	case "process-kill":
		killRandomProcess(req.KillProcess)
	default:
		log.Printf("[Agent] Unknown experiment type: %s", req.ExperimentType)
		return
	}

	log.Printf("[Agent] Experiment %s completed.", req.ExperimentType)
}

// -------------------------
// Network Fault Functions
// -------------------------

// addNetworkLatency adds latency to the default network interface via tc netem.
// This typically requires root privileges or a privileged container.
func addNetworkLatency(delayMs int) {
	if delayMs <= 0 {
		delayMs = 100
	}
	cmd := exec.Command("tc", "qdisc", "add", "dev", "eth0", "root", "netem", "delay", fmt.Sprintf("%dms", delayMs))
	if err := cmd.Run(); err != nil {
		log.Printf("[Agent] Failed to add network latency: %v", err)
	} else {
		log.Printf("[Agent] Added %dms network latency.", delayMs)
	}
}

// addNetworkLoss adds packet loss via tc netem.
func addNetworkLoss(lossPercent int) {
	if lossPercent <= 0 {
		lossPercent = 10
	}
	cmd := exec.Command("tc", "qdisc", "add", "dev", "eth0", "root", "netem", "loss", fmt.Sprintf("%d%%", lossPercent))
	if err := cmd.Run(); err != nil {
		log.Printf("[Agent] Failed to add network loss: %v", err)
	} else {
		log.Printf("[Agent] Added %d%% network loss.", lossPercent)
	}
}

// removeNetworkLatency removes the netem qdisc from eth0.
func removeNetworkLatency() {
	cmd := exec.Command("tc", "qdisc", "del", "dev", "eth0", "root")
	if err := cmd.Run(); err != nil {
		log.Printf("[Agent] Failed to remove netem qdisc: %v", err)
	} else {
		log.Println("[Agent] Removed network fault configuration.")
	}
}

// -------------------------
// Resource Starvation
// -------------------------

// runStressNGCPU runs stress-ng to stress CPU.
func runStressNGCPU(cpuWorkers, duration int) {
	if cpuWorkers <= 0 {
		cpuWorkers = 2
	}
	if duration <= 0 {
		duration = 30
	}
	args := []string{
		"--cpu", strconv.Itoa(cpuWorkers),
		"--timeout", fmt.Sprintf("%ds", duration),
	}
	log.Printf("[Agent] Running stress-ng CPU: %v", args)
	cmd := exec.Command("stress-ng", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[Agent] CPU stress-ng error: %v, output: %s", err, string(output))
	} else {
		log.Printf("[Agent] CPU stress-ng completed: %s", string(output))
	}
}

// runStressNGMemory runs stress-ng to stress memory.
func runStressNGMemory(memSizeMB, duration int) {
	if memSizeMB <= 0 {
		memSizeMB = 100
	}
	if duration <= 0 {
		duration = 30
	}
	args := []string{
		"--vm", "1",
		"--vm-bytes", fmt.Sprintf("%dm", memSizeMB),
		"--timeout", fmt.Sprintf("%ds", duration),
	}
	log.Printf("[Agent] Running stress-ng Memory: %v", args)
	cmd := exec.Command("stress-ng", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[Agent] Memory stress-ng error: %v, output: %s", err, string(output))
	} else {
		log.Printf("[Agent] Memory stress-ng completed: %s", string(output))
	}
}

// -------------------------
// Process Kill
// -------------------------

// killRandomProcess kills a random process or a specific process matching a pattern.
func killRandomProcess(pattern string) {
	log.Printf("[Agent] Attempting to kill a process matching pattern: %s", pattern)
	if pattern == "" {
		pattern = "go" // default pattern to kill something, but be careful!
	}

	// List all processes
	psCmd := exec.Command("ps", "-eo", "pid,cmd")
	out, err := psCmd.Output()
	if err != nil {
		log.Printf("[Agent] Failed to list processes: %v", err)
		return
	}

	lines := strings.Split(string(out), "\n")
	var candidates []string
	for _, line := range lines {
		if strings.Contains(line, pattern) {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				candidates = append(candidates, fields[0]) // pid is first
			}
		}
	}

	if len(candidates) == 0 {
		log.Printf("[Agent] No matching process found for pattern: %s", pattern)
		return
	}

	// Pick a random PID from the list
	rand.Seed(time.Now().UnixNano())
	pid := candidates[rand.Intn(len(candidates))]
	killCmd := exec.Command("kill", "-9", pid)
	if err := killCmd.Run(); err != nil {
		log.Printf("[Agent] Failed to kill process %s: %v", pid, err)
	} else {
		log.Printf("[Agent] Killed process %s matching pattern: %s", pid, pattern)
	}
}
