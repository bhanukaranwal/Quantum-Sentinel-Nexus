// services/digital-twin-svc/build.sbt

// --- Project Metadata ---
ThisBuild / version := "1.0.0"
ThisBuild / scalaVersion := "3.3.3" // Specify Scala 3

lazy val root = (project in file("."))
  .settings(
    name := "digital-twin-svc"
  )

// --- Akka Version ---
val AkkaVersion = "2.9.3"
val AkkaHttpVersion = "10.6.2"

// --- Dependencies ---
libraryDependencies ++= Seq(
  // Akka Streams for building reactive data processing pipelines
  "com.typesafe.akka" %% "akka-stream" % AkkaVersion,

  // Akka HTTP for building the HTTP/gRPC server
  "com.typesafe.akka" %% "akka-http" % AkkaHttpVersion,

  // Akka's gRPC library for defining and serving gRPC services
  "com.typesafe.akka" %% "akka-http-grpc" % AkkaHttpVersion,

  // SLF4J and Logback for structured logging
  "ch.qos.logback" % "logback-classic" % "1.5.6",
  "com.typesafe.akka" %% "akka-slf4j" % AkkaVersion,

  // Akka Testkit for testing
  "com.typesafe.akka" %% "akka-stream-testkit" % AkkaVersion % Test
)

// --- sbt-revolver ---
// A helpful plugin for running the application in development with hot-reloading.
addSbtPlugin("io.spray" % "sbt-revolver" % "0.10.0")

// --- Akka gRPC Codegen ---
// Enable the Akka gRPC plugin to generate code from .proto files.
enablePlugins(AkkaGrpcPlugin)
