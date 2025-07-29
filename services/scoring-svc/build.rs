// services/scoring-svc/build.rs

fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("cargo:rerun-if-changed=proto/scoring/v1/scoring.proto");

    // Use tonic_build to compile the protobuf file.
    // This will generate the necessary Rust code for the gRPC service.
    tonic_build::configure()
        // We specify the path to the directory containing our proto file.
        .compile(&["proto/scoring/v1/scoring.proto"], &["proto"])?;

    Ok(())
}
