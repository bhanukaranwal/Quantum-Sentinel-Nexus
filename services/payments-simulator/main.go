// services/payments-simulator/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

// --- Configuration ---
const (
	DefaultLocustHost = "http://ingestion-svc:8080"
	DefaultSpawnRate  = "100" // users to spawn per second
)

// --- HTTP Handlers ---

// handleStartSimulation triggers the Locust load test.
func handleStartSimulation(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	host := r.URL.Query().Get("host")
	if host == "" {
		host = DefaultLocustHost
	}

	users := r.URL.Query().Get("users")
	if users == "" {
		users = "1000" // Default to 1000 concurrent users
	}

	spawnRate := r.URL.Query().Get("spawn_rate")
	if spawnRate == "" {
		spawnRate = DefaultSpawnRate
	}

	runTime := r.URL.Query().Get("time")
	if runTime == "" {
		runTime = "60s" // Default to a 60-second run
	}

	log.Printf("🚀 Received request to start simulation with users=%s, spawn-rate=%s, time=%s, host=%s\n", users, spawnRate, runTime, host)

	// Build the Locust command
	cmd := exec.Command(
		"locust",
		"-f", "locustfile.py",
		"--headless", // Run without the web UI
		"-u", users,
		"-r", spawnRate,
		"-t", runTime,
		"--host", host,
	)

	// Set the command's output to our logger for visibility
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command asynchronously
	go func() {
		log.Println("🔥 Starting Locust load test...")
		if err := cmd.Run(); err != nil {
			log.Printf("❌ Locust command failed: %v\n", err)
		} else {
			log.Println("✅ Locust load test finished successfully.")
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Load test started. Check service logs for details.\n")
}

// handleHealthCheck provides a simple health endpoint.
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "ok"}`)
}

// --- Main Function ---

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/start", handleStartSimulation)
	mux.HandleFunc("/healthz", handleHealthCheck)

	port := 8089
	log.Printf("🌱 Payments Simulator listening on port %d\n", port)
	log.Printf("   -> Send requests to http://localhost:%d/start to begin a simulation.\n", port)

	if err := http.ListenAndServe(":"+strconv.Itoa(port), mux); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
