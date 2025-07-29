# .devcontainer/Dockerfile
#
# This Dockerfile defines the development environment for the Quantum Sentinel Nexus.
# It installs all the necessary languages, SDKs, and command-line tools.

# Start from a base Ubuntu image provided by Microsoft for dev containers.
FROM mcr.microsoft.com/devcontainers/base:ubuntu-22.04

# --- Install OS Packages & Tools ---
# Install common utilities and dependencies needed for various tools.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    curl \
    git \
    gnupg \
    lsb-release \
    software-properties-common \
    unzip \
    zip \
    # Cleanup
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# --- Install Language Runtimes ---

# Install Go
RUN go install golang.org/dl/go1.22.5@latest && go1.22.5 download

# Install Rust
ENV RUSTUP_HOME=/usr/local/rustup \
    CARGO_HOME=/usr/local/cargo \
    PATH=/usr/local/cargo/bin:$PATH
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

# Install Node.js (using nvm for flexibility)
ENV NVM_DIR="/usr/local/share/nvm"
ENV NODE_VERSION="18.18.2"
RUN mkdir -p $NVM_DIR && curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash \
    && . $NVM_DIR/nvm.sh \
    && nvm install $NODE_VERSION \
    && nvm use $NODE_VERSION \
    && nvm alias default $NODE_VERSION
ENV PATH $NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

# Install Python
RUN apt-get update && apt-get install -y python3.12 python3-pip \
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# --- Install Cloud & K8s Tooling ---

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Install Helm
RUN curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Install Kind
RUN go install sigs.k8s.io/kind@v0.23.0

# Install Tilt
RUN curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash

# Install Terraform
RUN curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add - \
    && apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main" \
    && apt-get update && apt-get install terraform

# --- Install qsnctl (Developer CLI) ---
# In a real project, you would download a pre-compiled binary from a release page.
# For now, we'll build it from source.
COPY services/tools/cli /tmp/qsnctl-src
RUN cd /tmp/qsnctl-src && cargo build --release \
    && mv target/release/qsnctl /usr/local/bin/

# --- Final Configuration ---
# Set a default user. The devcontainer.json will use this user.
USER vscode
