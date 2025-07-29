// services/compliance-svc/src/main/kotlin/com/qsn/kafka/AlertListener.kt
package com.qsn.kafka

import com.fasterxml.jackson.databind.ObjectMapper
import com.qsn.service.ComplianceService
import org.slf4j.LoggerFactory
import org.springframework.kafka.annotation.KafkaListener
import org.springframework.stereotype.Component

/**
 * A data class to represent the structure of an incoming alert message.
 */
data class Alert(
    val transactionId: String,
    val severity: String,
    val score: Double,
    val bankId: String,
    val modelVersion: String,
    val timestamp: String
)

/**
 * Listens for messages on the 'alerts' Kafka topic.
 */
@Component
class AlertListener(
    private val complianceService: ComplianceService,
    private val objectMapper: ObjectMapper
) {
    private val logger = LoggerFactory.getLogger(javaClass)

    /**
     * Consumes messages from the configured Kafka topic.
     * @param message The raw message payload as a String (JSON).
     */
    @KafkaListener(topics = ["\${kafka.topic.alerts}"], groupId = "\${kafka.group.id}")
    fun listen(message: String) {
        logger.info("📨 Received alert message: $message")
        try {
            // Deserialize the JSON string into our Alert data class
            val alert = objectMapper.readValue(message, Alert::class.java)

            // Pass the deserialized alert to the compliance service for processing
            complianceService.evaluateAlert(alert)
        } catch (e: Exception) {
            logger.error("❌ Failed to parse alert message: $message", e)
            // In a real application, you would send this message to a dead-letter queue.
        }
    }
}
