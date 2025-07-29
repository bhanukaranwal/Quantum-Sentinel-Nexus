# services/feature-engine-svc/main.py
import json
import os
import sys
import logging
from time import sleep

# Use the confluent-kafka library for robust Kafka integration.
from confluent_kafka import Consumer, Producer, KafkaError

# --- Logging Setup ---
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    stream=sys.stdout
)
logger = logging.getLogger(__name__)

# --- Configuration ---
KAFKA_BROKERS = os.environ.get('KAFKA_BROKERS', 'kafka:9092')
CONSUMER_GROUP_ID = os.environ.get('CONSUMER_GROUP_ID', 'feature-engine-group')
INPUT_TOPIC = os.environ.get('INPUT_TOPIC', 'txn-raw')
OUTPUT_TOPIC = os.environ.get('OUTPUT_TOPIC', 'features')


def create_kafka_consumer(brokers, group_id):
    """Creates and returns a Kafka Consumer."""
    conf = {
        'bootstrap.servers': brokers,
        'group.id': group_id,
        'auto.offset.reset': 'earliest', # Start reading from the beginning of the topic
        'enable.auto.commit': False # We will commit offsets manually
    }
    try:
        consumer = Consumer(conf)
        logger.info(f"✅ Kafka Consumer created successfully for brokers: {brokers}")
        return consumer
    except Exception as e:
        logger.error(f"❌ Failed to create Kafka Consumer: {e}")
        sys.exit(1)

def create_kafka_producer(brokers):
    """Creates and returns a Kafka Producer."""
    conf = {'bootstrap.servers': brokers}
    try:
        producer = Producer(conf)
        logger.info(f"✅ Kafka Producer created successfully for brokers: {brokers}")
        return producer
    except Exception as e:
        logger.error(f"❌ Failed to create Kafka Producer: {e}")
        sys.exit(1)

def feature_engineering(raw_txn: dict) -> dict:
    """
    Performs feature engineering on a raw transaction.
    In a real-world scenario, this would involve more complex logic, such as:
    - Looking up historical data from a database (e.g., Redis, Cassandra).
    - Using libraries like Pandas or RAPIDS for complex calculations.
    - Calculating features like:
        - Transaction frequency for the account.
        - Deviation from the account's average transaction amount.
        - Time of day/day of week features.
        - Whether the 'to_account' is new or has been seen before.
    """
    logger.info(f"Engineering features for transaction ID: {raw_txn.get('id')}")

    # For this example, we'll add a few simple features.
    features = {
        'transaction_id': raw_txn.get('id'),
        'original_amount': raw_txn.get('amount'),
        'currency': raw_txn.get('currency'),
        'bank_id': raw_txn.get('bank_id'),
        'is_high_value': raw_txn.get('amount', 0) > 10000,
        'is_international': not raw_txn.get('from_account', '').startswith(raw_txn.get('to_account', '')[:3]),
        'timestamp_utc': raw_txn.get('timestamp'),
    }
    
    # This is where you would integrate with RAPIDS for GPU-accelerated data processing.
    # For example:
    # gdf = cudf.DataFrame([raw_txn])
    # gdf['new_feature'] = gdf['amount'] * 2
    # features.update(gdf.to_pandas().to_dict('records')[0])

    return features

def delivery_report(err, msg):
    """
    Callback function for Kafka producer messages.
    Called once for each message produced to indicate delivery result.
    """
    if err is not None:
        logger.error(f"❌ Message delivery failed: {err}")
    else:
        logger.info(f"✅ Message delivered to {msg.topic()} [{msg.partition()}] at offset {msg.offset()}")

def main():
    """Main processing loop."""
    logger.info("🚀 Starting Feature Engine Service...")
    
    consumer = create_kafka_consumer(KAFKA_BROKERS, CONSUMER_GROUP_ID)
    producer = create_kafka_producer(KAFKA_BROKERS)

    consumer.subscribe([INPUT_TOPIC])
    logger.info(f"👂 Subscribed to topic: {INPUT_TOPIC}")

    try:
        while True:
            msg = consumer.poll(timeout=1.0) # Poll for new messages

            if msg is None:
                continue
            if msg.error():
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    # End of partition event
                    logger.info(f"Reached end of partition for {msg.topic()} [{msg.partition()}]")
                else:
                    logger.error(f"❌ Kafka consumer error: {msg.error()}")
                    # Consider a backoff strategy here
                    sleep(5)
                continue

            # Process the message
            try:
                raw_txn = json.loads(msg.value().decode('utf-8'))
                logger.info(f"📨 Consumed message from {msg.topic()}: key={msg.key()} value={raw_txn}")

                # Perform feature engineering
                enriched_features = feature_engineering(raw_txn)

                # Produce the enriched features to the output topic
                producer.produce(
                    OUTPUT_TOPIC,
                    key=str(enriched_features.get('transaction_id')),
                    value=json.dumps(enriched_features).encode('utf-8'),
                    callback=delivery_report
                )
                
                # Manually commit the offset after successful processing.
                consumer.commit(asynchronous=True)

            except json.JSONDecodeError:
                logger.error(f"Malformed JSON received: {msg.value().decode('utf-8')}")
                # Decide on a dead-letter queue strategy for bad messages
            except Exception as e:
                logger.error(f"An unexpected error occurred: {e}")

            # The producer is asynchronous, poll to trigger callbacks
            producer.poll(0)

    except KeyboardInterrupt:
        logger.info("🔌 User interrupted, shutting down...")
    finally:
        # Cleanly close the consumer and producer
        logger.info("Closing Kafka consumer...")
        consumer.close()
        logger.info("Flushing and closing Kafka producer...")
        producer.flush()


if __name__ == '__main__':
    main()

