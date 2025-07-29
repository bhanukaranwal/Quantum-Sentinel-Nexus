// services/compliance-svc/src/main/kotlin/com/qsn/configuration/DroolsConfiguration.kt
package com.qsn.configuration

import org.kie.api.KieServices
import org.kie.api.runtime.KieContainer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

@Configuration
class DroolsConfiguration {

    /**
     * Creates the KieContainer bean, which is the central point for Drools rule engine access.
     * It loads the rules defined in the project's resources.
     * @return A KieContainer instance.
     */
    @Bean
    fun kieContainer(): KieContainer {
        val kieServices = KieServices.Factory.get()
        return kieServices.kieClasspathContainer
    }
}
