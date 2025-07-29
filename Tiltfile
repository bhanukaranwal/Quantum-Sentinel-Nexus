# Tiltfile
# This file configures the live development environment for the Quantum Sentinel Nexus.
# It tells Tilt how to build container images, deploy services to Kubernetes,
# and set up live-reloading for a fast development loop.
#
# To run: `tilt up` or `qsnctl dev-up`

# --- Configuration ---

# Load extensions for common tasks like deploying Helm charts.
load('ext://helm_resource', 'helm_resource')
load('ext://kubectl', 'kubectl')

# Define a list of all microservices to be managed by Tilt.
# This makes it easy to add or remove services.
SERVICES = [
    "ingestion-svc",
    "feature-engine-svc",
    "scoring-svc",
    "fl-coordinator-svc",
    "alert-router-svc",
    "audit-svc", # Note: Chaincode deployment is more complex, this is a placeholder
    "auth-gateway-svc",
    "dashboard-ui",
    "threat-intel-svc",
    "digital-twin-svc",
    "edge-node-agent", # Placeholder for edge simulation
    "compliance-svc",
    "neuromorph-svc", # Placeholder for neuromorphic simulation
    "payments-simulator",
]

# --- Global Setup ---

# First, deploy the core infrastructure dependencies using their Helm charts.
# This includes Kafka, Cassandra, Neo4j, etc.
# In a real project, you would point to your custom charts or a community chart.
helm_resource(
    name='infra-dependencies',
    chart='charts/qsn-infra',
    release_name='qsn-infra',
    namespace='default',
    wait=True,
    wait_timeout='5m'
)

# --- Service Definitions ---

# Loop through the services list and create Tilt resources for each one.
# This avoids repetitive code and makes the file easier to maintain.
for svc_name in SERVICES:
    # Define the Docker build for the service.
    # The context is the service's directory.
    # The Dockerfile path is standard.
    docker_build(
        f'qsn/{svc_name}',
        f'services/{svc_name}',
        dockerfile=f'services/{svc_name}/Dockerfile',
        # For services with compiled code (Go, Rust), live_update is more complex.
        # For interpreted languages (Python, Node) or frontends, we can sync files.
        live_update=[
            # Sync local file changes into the running container to avoid a full rebuild.
            # This is highly effective for Node.js, Python, and UI development.
            # Example for a Node.js service:
            # sync(f'services/{svc_name}/src', f'/app/src'),
            # run('npm install', trigger=f'services/{svc_name}/package.json'),
            # run('npm run restart', trigger=f'services/{svc_name}/src')
        ] if svc_name in ['dashboard-ui', 'auth-gateway-svc', 'fl-coordinator-svc'] else []
    )

    # Define the Kubernetes deployment for the service.
    # We use a Helm chart for each service, which is best practice.
    # The `k8s_resource` function tells Tilt to deploy and manage this resource.
    helm_resource(
        name=svc_name,
        chart=f'charts/{svc_name}',
        release_name=svc_name,
        namespace='default',
        # This tells Tilt to wait for the deployment to be ready before proceeding.
        wait=True,
        # This links the k8s resource to the Docker image we defined above.
        # When the image is rebuilt, Tilt knows to update this deployment.
        image_deps=[f'qsn/{svc_name}'],
        # Set up port forwarding to access the service from localhost.
        port_forwards=[
            # Example: '3000:3000' for the dashboard-ui
        ]
    )

# --- Custom Resource Groups ---

# Group resources in the Tilt UI for better organization.
ui.resource_group(
    name='core-backend',
    resources=[
        'ingestion-svc',
        'feature-engine-svc',
        'scoring-svc',
        'fl-coordinator-svc',
    ]
)

ui.resource_group(
    name='support-services',
    resources=[
        'alert-router-svc',
        'audit-svc',
        'auth-gateway-svc',
        'compliance-svc',
    ]
)

ui.resource_group(
    name='ui-and-simulators',
    resources=[
        'dashboard-ui',
        'payments-simulator',
    ]
)

# --- Final Message ---
print("Tiltfile loaded. Starting services for Quantum Sentinel Nexus...")
print("Navigate to http://localhost:10350 to view the Tilt UI.")

