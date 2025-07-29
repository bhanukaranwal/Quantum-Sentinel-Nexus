// services/compliance-svc/src/main/kotlin/com/qsn/ComplianceSvcApplication.kt
package com.qsn

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

/**
 * The main entry point for the Compliance Service application.
 * The @SpringBootApplication annotation is a convenience annotation that adds all of the following:
 * - @Configuration: Tags the class as a source of bean definitions for the application context.
 * - @EnableAutoConfiguration: Tells Spring Boot to start adding beans based on classpath settings, other beans, and various property settings.
 * - @ComponentScan: Tells Spring to look for other components, configurations, and services in the 'com.qsn' package, allowing it to find and register them.
 */
@SpringBootApplication
class ComplianceSvcApplication

/**
 * The main function that runs the Spring Boot application.
 */
fun main(args: Array<String>) {
    runApplication<ComplianceSvcApplication>(*args)
}
