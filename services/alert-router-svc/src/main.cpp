// services/alert-router-svc/src/main.cpp
#include <iostream>
#include <string>
#include <thread>
#include <vector>
#include <csignal>
#include <atomic>

// High-performance WebSocket library
#include "uWebSockets/App.h"
// C++ Kafka client library
#include "librdkafka/rdkafkacpp.h"
// JSON library for parsing messages
#include "nlohmann/json.hpp"
// ZeroMQ library for fallback/internal messaging
#include "zmq.hpp"

// Use namespaces for convenience
using json = nlohmann::json;

// Global flag to handle graceful shutdown
std::atomic<bool> running = true;

void signal_handler(int signum) {
    std::cout << "\nInterrupt signal (" << signum << ") received.\n";
    running = false;
}

// --- Kafka Consumer Logic ---

class KafkaConsumer {
public:
    KafkaConsumer(const std::string& brokers, const std::string& group_id, uWS::App* ws_app)
        : brokers_(brokers), group_id_(group_id), ws_app_(ws_app) {
        
        // Create Kafka configuration
        conf_ = RdKafka::Conf::create(RdKafka::Conf::CONF_GLOBAL);
        RdKafka::Conf::create(RdKafka::Conf::CONF_TOPIC);

        std::string errstr;
        conf_->set("bootstrap.servers", brokers_, errstr);
        conf_->set("group.id", group_id_, errstr);
        conf_->set("auto.offset.reset", "earliest", errstr);

        // Create consumer
        consumer_ = RdKafka::KafkaConsumer::create(conf_, errstr);
        if (!consumer_) {
            std::cerr << "Failed to create consumer: " << errstr << std::endl;
            exit(1);
        }
        std::cout << "✅ Kafka consumer created" << std::endl;
    }

    ~KafkaConsumer() {
        consumer_->close();
        delete consumer_;
        delete conf_;
    }

    void consume(const std::string& topic) {
        std::vector<std::string> topics = {topic};
        consumer_->subscribe(topics);
        std::cout << "👂 Subscribed to topic: " << topic << std::endl;

        while (running) {
            RdKafka::Message *msg = consumer_->consume(1000); // 1 second timeout
            if (!msg) continue;

            handle_message(msg);
            delete msg;
        }
    }

private:
    void handle_message(RdKafka::Message* message) {
        switch (message->err()) {
            case RdKafka::ERR_NO_ERROR: {
                // We have a valid message
                std::string payload = static_cast<const char *>(message->payload());
                std::cout << "📨 Consumed message from Kafka: " << payload << std::endl;

                // TODO: Implement routing logic based on rules engine (e.g., Drools)
                // For now, we broadcast all alerts to all WebSocket clients.
                try {
                    json alert = json::parse(payload);
                    std::string alert_type = alert.value("severity", "INFO"); // e.g., "CRITICAL", "HIGH"
                    
                    // The topic for WebSocket broadcast could be based on severity, bank_id, etc.
                    std::string ws_topic = "alerts/" + alert_type;

                    // Publish to the WebSocket server on the main thread
                    ws_app_->publish(ws_topic, payload, uWS::OpCode::TEXT);
                } catch (const json::parse_error& e) {
                    std::cerr << "JSON parse error: " << e.what() << std::endl;
                }
                break;
            }
            case RdKafka::ERR__PARTITION_EOF:
                // End of partition, not an error
                break;
            default:
                std::cerr << "Consume error: " << message->errstr() << std::endl;
                break;
        }
    }

    std::string brokers_;
    std::string group_id_;
    RdKafka::Conf *conf_;
    RdKafka::KafkaConsumer *consumer_;
    uWS::App* ws_app_; // Pointer to the WebSocket App to publish messages
};

// --- Main Application ---

int main() {
    // Register signal handler for graceful shutdown
    signal(SIGINT, signal_handler);
    signal(SIGTERM, signal_handler);

    // --- Configuration ---
    std::string kafka_brokers = std::getenv("KAFKA_BROKERS") ? std::getenv("KAFKA_BROKERS") : "kafka:9092";
    std::string kafka_topic = std::getenv("KAFKA_TOPIC") ? std::getenv("KAFKA_TOPIC") : "alerts";
    std::string kafka_group_id = "alert-router-group";
    int ws_port = 8081;

    // --- WebSocket Server Setup ---
    uWS::App ws_app = uWS::App().ws<std::string>("/*", {
        /* Settings */
        .compression = uWS::SHARED_COMPRESSOR,
        .maxPayloadLength = 16 * 1024,
        /* Handlers */
        .open = [](auto *ws) {
            std::cout << "WebSocket client connected." << std::endl;
            // Subscribe to all alert types by default
            ws->subscribe("alerts/CRITICAL");
            ws->subscribe("alerts/HIGH");
            ws->subscribe("alerts/INFO");
        },
        .message = [](auto *ws, std::string_view message, uWS::OpCode opCode) {
            // Handle incoming messages from clients, e.g., for changing subscriptions
            std::cout << "Received message from client: " << message << std::endl;
        },
        .close = [](auto *ws, int code, std::string_view message) {
            std::cout << "WebSocket client disconnected." << std::endl;
        }
    }).listen(ws_port, [ws_port](auto *listen_socket) {
        if (listen_socket) {
            std::cout << "✅ WebSocket server listening on port " << ws_port << std::endl;
        }
    });

    // --- ZeroMQ Fallback/Internal Pub-Sub Setup ---
    // This could be used for inter-process communication or as a fallback if Kafka is down.
    zmq::context_t zmq_context(1);
    zmq::socket_t zmq_publisher(zmq_context, zmq::socket_type::pub);
    zmq_publisher.bind("tcp://*:5555");
    std::cout << "✅ ZeroMQ publisher bound to tcp://*:5555" << std::endl;

    // --- Start Kafka Consumer Thread ---
    KafkaConsumer kafka_consumer(kafka_brokers, kafka_group_id, &ws_app);
    std::thread kafka_thread(&KafkaConsumer::consume, &kafka_consumer, kafka_topic);
    
    std::cout << "🚀 Alert Router Service started." << std::endl;
    
    // Run the WebSocket server event loop in the main thread
    ws_app.run();

    // --- Shutdown ---
    std::cout << "WebSocket server stopped. Waiting for Kafka thread to join..." << std::endl;
    kafka_thread.join();

    std::cout << "✅ Service shut down gracefully." << std::endl;

    return 0;
}
