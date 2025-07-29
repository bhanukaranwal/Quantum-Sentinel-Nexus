// services/ingestion-svc/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	// Uncomment the following line to enable Kafka integration
	// "github.com/twmb/franz-go/pkg/kgo"
)

// --- Constants ---

const (
	KafkaBrokerEnv = "KAFKA_BROKERS"
	KafkaTopicEnv  = "KAFKA_TOPIC"
	DefaultBroker  = "kafka:9092"
	DefaultTopic   = "txn-raw"
)

// --- Structs ---

// Transaction represents the structure of an incoming transaction request.
// In a real application, this would be much more detailed.
type Transaction struct {
	ID            string    `json:"id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	FromAccount   string    `json:"from_account"`
	ToAccount     string    `json:"to_account"`
	Timestamp     time.Time `json:"timestamp"`
	BankID        string    `json:"bank_id"`
	ClientDetails string    `json:"client_details"` // e.g., IP, User-Agent
}

// IngestionService holds the service's dependencies, like the Kafka producer.
type IngestionService struct {
	// Uncomment the following line to enable Kafka integration
	// kafkaProducer *kgo.Client
	router *chi.Mux
}

// --- Main Function ---

func main() {
	log.Println("🚀 Starting Ingestion Service...")

	// Get configuration from environment variables
	kafkaBrokers := getEnv(KafkaBrokerEnv, DefaultBroker)
	kafkaTopic := getEnv(KafkaTopicEnv, DefaultTopic)

	// --- Dependency Injection ---
	// In a real app, you might use a framework like Wire or fx.

	// Uncomment the following lines to create the Kafka producer
	/*
		seeds := []string{kafkaBrokers}
		producer, err := kgo.NewClient(
			kgo.SeedBrokers(seeds...),
			kgo.DefaultProduceTopic(kafkaTopic),
			kgo.AllowAutoTopicCreation(),
		)
		if err != nil {
			log.Fatalf("❌ Failed to create Kafka client: %v", err)
		}
		defer producer.Close()
		log.Printf("✅ Connected to Kafka broker at %s", kafkaBrokers)
	*/

	service := &IngestionService{
		// kafkaProducer: producer,
		router: chi.NewRouter(),
	}

	// --- Middleware ---
	service.router.Use(middleware.RequestID)
	service.router.Use(middleware.RealIP)
	service.router.Use(middleware.Logger)
	service.router.Use(middleware.Recoverer)
	service.router.Use(middleware.Timeout(60 * time.Second))

	// --- Routes ---
	service.router.Post("/v1/txn", service.handleTransactionIngest(kafkaTopic))
	service.router.Get("/healthz", handleHealthCheck)

	// --- Server Startup & Graceful Shutdown ---
	server := &http.Server{
		Addr:    ":8080",
		Handler: service.router,
	}

	go func() {
		log.Println("👂 Listening on port :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8080: %v\n", err)
		}
	}()

	// Wait for termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("🔌 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %+v", err)
	}

	log.Println("✅ Server gracefully stopped")
}

// --- HTTP Handlers ---

// handleTransactionIngest is the primary handler for receiving transactions.
func (s *IngestionService) handleTransactionIngest(topic string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var txn Transaction

		// 1. Decode JSON body
		if err := json.NewDecoder(r.Body).Decode(&txn); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 2. TODO: Implement Post-Quantum Signature Verification
		// This would involve reading a signature from the request headers (e.g., 'X-Dilithium-Signature')
		// and verifying it against the request payload using a PQC library.
		// pqc.VerifyDilithium(signature, payload, publicKey)

		// 3. Validate and enrich the transaction
		if txn.ID == "" {
			txn.ID = uuid.New().String()
		}
		if txn.Timestamp.IsZero() {
			txn.Timestamp = time.Now().UTC()
		}

		// 4. Serialize to Protobuf for Kafka
		// In a real implementation, you would convert the JSON struct to a Protobuf message
		// for efficiency and schema enforcement.
		payload, err := json.Marshal(txn)
		if err != nil {
			http.Error(w, "Failed to serialize transaction", http.StatusInternalServerError)
			return
		}

		// 5. Publish to Kafka
		// Uncomment the following block to enable Kafka publishing
		/*
			record := &kgo.Record{Topic: topic, Value: payload}
			s.kafkaProducer.Produce(r.Context(), record, func(r *kgo.Record, err error) {
				if err != nil {
					log.Printf("❌ Kafka produce error: %v\n", err)
					// Implement a dead-letter queue or other retry mechanism here.
				} else {
					log.Printf("✅ Produced record to topic %s, partition %d, offset %d\n", r.Topic, r.Partition, r.Offset)
				}
			})
		*/

		// For demonstration without Kafka, we just log it.
		log.Printf("Received transaction: %s", payload)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"status": "accepted", "transaction_id": txn.ID})
	}
}

// handleHealthCheck provides a simple health endpoint.
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// --- Helper Functions ---

// getEnv is a helper to read an environment variable or return a default value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
