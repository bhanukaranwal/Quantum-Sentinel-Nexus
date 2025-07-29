# services/threat-intel-svc/main.py
import os
import logging
import sys
import asyncio
from fastapi import FastAPI, HTTPException
from neo4j import AsyncGraphDatabase, exceptions
from pydantic import BaseModel
# from confluent_kafka import Consumer # Placeholder for Kafka integration

# --- Logging Setup ---
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    stream=sys.stdout
)
logger = logging.getLogger(__name__)

# --- Configuration ---
NEO4J_URI = os.environ.get("NEO4J_URI", "neo4j://neo4j:7687")
NEO4J_USER = os.environ.get("NEO4J_USER", "neo4j")
NEO4J_PASSWORD = os.environ.get("NEO4J_PASSWORD", "password") # Use a secure password in production
KAFKA_BROKERS = os.environ.get('KAFKA_BROKERS', 'kafka:9092')
KAFKA_TOPIC = os.environ.get('KAFKA_TOPIC', 'entity-events')

# --- FastAPI Application ---
app = FastAPI(
    title="Threat Intel Service",
    description="Provides graph-based threat intelligence from a Neo4j database.",
    version="1.0.0",
)

# --- Neo4j Driver ---
# The driver is a global instance that will be managed by the app's lifespan events.
driver = None

@app.on_event("startup")
async def startup_event():
    """Initializes the Neo4j driver on application startup."""
    global driver
    try:
        driver = AsyncGraphDatabase.driver(NEO4J_URI, auth=(NEO4J_USER, NEO4J_PASSWORD))
        await driver.verify_connectivity()
        logger.info("✅ Neo4j driver connected and verified.")
        # asyncio.create_task(kafka_consumer_task()) # Start Kafka consumer in the background
    except exceptions.ServiceUnavailable as e:
        logger.error(f"❌ Could not connect to Neo4j at {NEO4J_URI}: {e}")
        # In a real app, you might want to exit or have a retry mechanism.
        driver = None

@app.on_event("shutdown")
async def shutdown_event():
    """Closes the Neo4j driver on application shutdown."""
    if driver:
        await driver.close()
        logger.info("🔌 Neo4j driver closed.")


# --- Pydantic Models for API ---
class Entity(BaseModel):
    id: str
    type: str # e.g., "ACCOUNT", "DEVICE", "USER"

class Relationship(BaseModel):
    source: Entity
    target: Entity
    type: str # e.g., "HAS_USED", "SENT_FUNDS_TO"


# --- API Endpoints ---

@app.get("/healthz", summary="Health Check")
async def health_check():
    """Checks the health of the service and its connection to Neo4j."""
    if not driver:
        raise HTTPException(status_code=503, detail="Neo4j connection not available.")
    try:
        await driver.verify_connectivity()
        return {"status": "ok", "neo4j_connection": "ok"}
    except exceptions.ServiceUnavailable:
        raise HTTPException(status_code=503, detail="Neo4j connection is down.")

@app.get("/entity/{entity_type}/{entity_id}/related", summary="Find Related Entities")
async def get_related_entities(entity_type: str, entity_id: str, depth: int = 1):
    """
    Finds all entities related to a given entity up to a certain depth.
    """
    if not driver:
        raise HTTPException(status_code=503, detail="Neo4j connection not available.")

    query = f"""
    MATCH (n:{entity_type.upper()} {{id: $entity_id}})-[*1..{depth}]-(related)
    RETURN DISTINCT labels(related) as type, related.id as id
    """
    
    async with driver.session() as session:
        result = await session.run(query, entity_id=entity_id)
        related = [
            {"type": record["type"][0], "id": record["id"]}
            async for record in result
        ]
    return {"entity": {"type": entity_type, "id": entity_id}, "related": related}

# --- Background Kafka Consumer (Placeholder) ---

async def kafka_consumer_task():
    """
    A background task to consume entity events from Kafka and update the graph.
    """
    logger.info("Starting Kafka consumer task...")
    # consumer = Consumer(...)
    # consumer.subscribe([KAFKA_TOPIC])
    while True:
        # msg = consumer.poll(1.0)
        # if msg is None: continue
        # if msg.error(): ...
        
        # event = json.loads(msg.value())
        # await update_graph(event)
        await asyncio.sleep(1) # Placeholder sleep

async def update_graph(event: dict):
    """
    Updates the Neo4j graph based on an incoming event.
    This function would contain Cypher queries to create or merge nodes and relationships.
    """
    if not driver:
        logger.error("Cannot update graph, Neo4j driver is not available.")
        return

    # Example: Create a relationship based on an event
    # event = {"source": {"type": "USER", "id": "u1"}, "target": {"type": "DEVICE", "id": "d1"}, "type": "LOGGED_IN_WITH"}
    
    query = """
    MERGE (s:{source_type} {{id: $source_id}})
    MERGE (t:{target_type} {{id: $target_id}})
    MERGE (s)-[r:{rel_type}]->(t)
    RETURN type(r)
    """.format(
        source_type=event['source']['type'],
        target_type=event['target']['type'],
        rel_type=event['type']
    )
    
    params = {
        "source_id": event['source']['id'],
        "target_id": event['target']['id'],
    }

    async with driver.session() as session:
        await session.run(query, params)
        logger.info(f"Updated graph with event: {event}")

