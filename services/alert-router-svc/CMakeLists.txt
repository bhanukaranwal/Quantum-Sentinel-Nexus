# services/alert-router-svc/CMakeLists.txt

# Set the minimum required version of CMake.
cmake_minimum_required(VERSION 3.16)

# Define the project name and specify C++ as the language.
project(alert-router-svc CXX)

# Set the C++ standard to C++20, as requested in the project specification.
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# --- Find Dependencies ---
# Use find_package to locate the necessary libraries and their headers.
# This assumes the libraries have been installed in a standard location
# or that their paths are provided to CMake during the build process
# (e.g., inside the Docker container).

# Find the ZeroMQ library.
find_package(ZeroMQ REQUIRED)

# Find the librdkafka library.
find_package(RdKafka REQUIRED)

# Find the uWebSockets library.
# uWebSockets often requires finding its own dependencies (libuv, OpenSSL).
find_package(uWebSockets REQUIRED)

# --- Include Directories ---
# Add the directories containing the header files for our dependencies.
# This allows us to use #include <...> for these libraries in our source code.

# The nlohmann/json library is often used as a header-only library,
# so we just need to tell CMake where to find its 'include' directory.
# We assume it's in a 'vendor' directory for this example.
include_directories(vendor/json/include)

# Add include directories for the found packages.
include_directories(${ZMQ_INCLUDE_DIRS})
include_directories(${RDKAFKA_INCLUDE_DIRS})
include_directories(${UWEBSOCKETS_INCLUDE_DIRS})


# --- Executable Definition ---

# Define the name of the executable to be built.
add_executable(alert-router-svc src/main.cpp)

# --- Linking Libraries ---
# Link the executable with the libraries it depends on.
target_link_libraries(alert-router-svc
    PRIVATE
    # Link the actual library targets provided by find_package.
    RdKafka::rdkafkacpp
    uWebSockets::uWS
    ZeroMQ::ZeroMQ

    # We also need to link pthreads, which is used by several of these libraries.
    pthread
)

# --- Installation ---
# This section defines how to install the built executable.
# This is useful for packaging the application.
install(TARGETS alert-router-svc
    DESTINATION bin
)
