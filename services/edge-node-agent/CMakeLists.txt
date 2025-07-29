# services/edge-node-agent/CMakeLists.txt

# Set the minimum required version of CMake.
cmake_minimum_required(VERSION 3.16)

# Define the project name and specify C++ as the language.
project(edge-node-agent CXX)

# Set the C++ standard to C++20.
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# --- Find Dependencies ---
# Use find_package to locate the necessary libraries and their headers.
# This assumes the libraries are installed in a standard location.
# For cross-compilation (e.g., for NVIDIA Jetson), a toolchain file
# would be used to specify the paths to the target architecture's libraries.

# Find gRPC and Protobuf.
find_package(gRPC CONFIG REQUIRED)
find_package(Protobuf CONFIG REQUIRED)

# Find the CUDA toolkit for GPU acceleration.
# This is optional; the build will proceed without it, but GPU features
# will be disabled.
find_package(CUDA)
if(CUDA_FOUND)
    message(STATUS "CUDA Toolkit found. Enabling GPU features.")
    # This enables CUDA language support in the project.
    enable_language(CUDA)
else()
    message(STATUS "CUDA Toolkit not found. Building with CPU-only features.")
endif()


# --- Executable Definition ---

# Define the name of the executable to be built.
add_executable(edge-node-agent src/main.cpp)


# --- Linking Libraries ---
# Link the executable with the libraries it depends on.
target_link_libraries(edge-node-agent
    PRIVATE
    # Link the gRPC and Protobuf libraries.
    gRPC::grpc++
    protobuf::libprotobuf
)

# Add CUDA-specific linking if the toolkit was found.
if(CUDA_FOUND)
    target_link_libraries(edge-node-agent
        PRIVATE
        # Link the CUDA runtime library.
        CUDA::cudart
    )
endif()


# --- Installation ---
# This section defines how to install the built executable.
install(TARGETS edge-node-agent
    DESTINATION bin
)
