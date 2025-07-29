// services/compliance-svc/src/main/kotlin/com/qsn/service/ComplianceService.kt
package com.qsn.service

import com.qsn.kafka.Alert
import org.kie.api.runtime.KieContainer
import org.slf4j.LoggerFactory
import org.springframework.stereotype.Service

@Service
class ComplianceService(
    private val kieContainer: KieContainer
) {
    private val logger = LoggerFactory.getLogger(javaClass)

    /**
     * Evaluates a given alert against the configured Drools rules.
     * @param alert The alert object to be evaluated.
     */
    fun evaluateAlert(alert: Alert) {
        try {
            // Create a new session from the KieContainer for each request.
            // This is important for thread safety.
            val kieSession = kieContainer.newKieSession()
            logger.info("Created new KieSession for transaction: ${alert.transactionId}")

            // Insert the alert object into the session. The rules engine can now access it.
            kieSession.insert(alert)

            // Fire all rules that match the current state of the session.
            val rulesFired = kieSession.fireAllRules()
            logger.info("🔥 Fired $rulesFired rules for transaction: ${alert.transactionId}")

            // Dispose of the session to release its resources.
            kieSession.dispose()

        } catch (e: Exception) {
            logger.error("❌ Error during Drools rule evaluation for alert: $alert", e)
        }
    }
}
