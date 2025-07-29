# services/audit-svc/Dockerfile
#
# This Dockerfile is used to build the chaincode binary.
# The Hyperledger Fabric peer will use this environment to build the chaincode
# before launching it in its own container.

# We use the official Fabric-ccenv image, which contains all the necessary
# Go toolchains and libraries for building Fabric chaincode.
FROM hyperledger/fabric-ccenv:latest

# Set the working directory for the chaincode.
WORKDIR /chaincode

# Copy the chaincode source code into the build environment.
COPY . .

# The Fabric peer will execute this command to build the chaincode.
# The output binary will be placed in /chaincode/bin.
# This step is typically not run by a user directly, but by the peer's
# chaincode lifecycle process.
RUN go build -o /chaincode/bin/audit-svc

# The final chaincode container launched by the peer will use a different
# entrypoint. This Dockerfile is purely for defining the build environment.
