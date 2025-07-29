// services/digital-twin-svc/src/main/scala/com/qsn/DigitalTwinServer.scala
package com.qsn

import akka.actor.typed.ActorSystem
import akka.actor.typed.scaladsl.Behaviors
import akka.http.scaladsl.Http
import akka.http.scaladsl.model.{HttpRequest, HttpResponse}
import com.typesafe.config.ConfigFactory

import scala.concurrent.{ExecutionContext, Future}
import scala.concurrent.duration._
import scala.util.{Failure, Success}

// Import the generated gRPC classes
import com.qsn.grpc._

// --- Main Application Object ---
object DigitalTwinServer {

  def main(args: Array[String]): Unit = {
    // Load configuration from application.conf
    val conf = ConfigFactory
      .parseString("akka.http.server.preview.enable-http2 = on")
      .withFallback(ConfigFactory.defaultApplication())

    // Create the main ActorSystem, which is the heart of an Akka application
    implicit val system: ActorSystem[_] = ActorSystem(Behaviors.empty, "DigitalTwinService", conf)
    implicit val ec: ExecutionContext = system.executionContext

    // Create an instance of the gRPC service implementation
    val service: DigitalTwinService = new DigitalTwinServiceImpl()

    // Bind the gRPC service to a port
    val binding = Http()
      .newServerAt("0.0.0.0", 8082)
      .bind(DigitalTwinServiceHandler(service))

    binding.onComplete {
      case Success(binding) =>
        val address = binding.localAddress
        system.log.info(s"🚀 Digital Twin gRPC service bound to ${address.getHostString}:${address.getPort}")
      case Failure(ex) =>
        system.log.error("Failed to bind gRPC endpoint, terminating system", ex)
        system.terminate()
    }

    // TODO: Initialize and run the Akka Streams Kafka consumer pipeline here
    // DigitalTwinConsumer.run()
  }
}

// --- gRPC Service Implementation ---
class DigitalTwinServiceImpl(implicit system: ActorSystem[_]) extends DigitalTwinService {
  implicit val ec: ExecutionContext = system.executionContext
  private val log = system.log

  // This is a placeholder for a real, stateful digital twin store.
  // In a real Akka application, you would use Akka Actors, one for each customer,
  // to maintain state in memory. This state could be backed by a persistent store
  // like Cassandra.
  private val twinStore = Map(
    "customer-123" -> DigitalTwin(
      customerId = "customer-123",
      averageTxnValue = 150.75,
      txnFrequencyPerDay = 2.5,
      lastSeenLocation = "New York, NY"
    )
  )

  /**
   * gRPC method to retrieve a digital twin for a given customer.
   */
  override def getDigitalTwin(in: GetDigitalTwinRequest): Future[GetDigitalTwinResponse] = {
    log.info(s"Received request for digital twin of customer: ${in.customerId}")

    twinStore.get(in.customerId) match {
      case Some(twin) =>
        // Found the twin, return it in the response
        Future.successful(GetDigitalTwinResponse(Some(twin)))
      case None =>
        // No twin found for this customer
        log.warn(s"No digital twin found for customer: ${in.customerId}")
        Future.successful(GetDigitalTwinResponse(None))
    }
  }

  /**
    * gRPC method to simulate a transaction against a digital twin.
    */
  override def simulateTransaction(in: SimulateTransactionRequest): Future[SimulateTransactionResponse] = {
    log.info(s"Simulating transaction for customer: ${in.customerId} with amount ${in.amount}")

    // In a real system, this would involve sending a message to the specific
    // customer's actor, which would then update its internal state and run
    // simulation logic.
    val isAnomaly = in.amount > 5000 // Simple anomaly rule
    
    Future.successful(
      SimulateTransactionResponse(
        customerId = in.customerId,
        isAnomaly = isAnomaly,
        comment = if (isAnomaly) "Transaction value is significantly higher than average." else "Transaction appears normal."
      )
    )
  }
}
