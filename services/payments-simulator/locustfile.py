# services/payments-simulator/locustfile.py
import random
import uuid
from datetime import datetime
from locust import HttpUser, task, between

# --- Test Data ---
# In a real test, this data would be more extensive and realistic.
CURRENCIES = ["USD", "EUR", "GBP", "JPY"]
BANK_IDS = ["bank-a", "bank-b", "bank-c", "bank-d"]

def generate_account_id():
    """Generates a random account ID."""
    return f"acc_{random.randint(100000, 999999)}"

# --- Locust User Behavior ---

class TransactionUser(HttpUser):
    """
    Defines the behavior of a single user (or "locust").
    Each user will repeatedly execute the tasks defined below.
    """
    # The 'wait_time' defines how long a user will wait between tasks.
    # 'between(0.1, 0.5)' means a random wait time between 100 and 500 milliseconds.
    wait_time = between(0.1, 0.5)

    @task
    def send_transaction(self):
        """
        This is the main task for the user. It simulates sending a single
        fraudulent or legitimate transaction to the ingestion service.
        """
        headers = {
            "Content-Type": "application/json",
            # In a real system, you would add a post-quantum signature here.
            # "X-Dilithium-Signature": "signed_payload_hash"
        }
        
        # Create a synthetic transaction payload.
        payload = {
            "id": str(uuid.uuid4()),
            "amount": round(random.uniform(5.0, 5000.0), 2),
            "currency": random.choice(CURRENCIES),
            "from_account": generate_account_id(),
            "to_account": generate_account_id(),
            "timestamp": datetime.utcnow().isoformat() + "Z",
            "bank_id": random.choice(BANK_IDS),
        }

        # Send an HTTP POST request to the /v1/txn endpoint of the target host.
        # The target host is specified when starting the Locust test (e.g., http://ingestion-svc:8080).
        self.client.post("/v1/txn", json=payload, headers=headers)

