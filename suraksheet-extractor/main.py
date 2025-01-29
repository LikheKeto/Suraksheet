import os
import re
import signal
import pika
import json
import psycopg2
from minio import Minio
import pytesseract
from PIL import Image
import logging
import cv2
import numpy as np
from elasticsearch import Elasticsearch

es = Elasticsearch(os.getenv("ELASTICSEARCH_URL"))


def clean_text(text):
    # Remove unwanted characters and extra spaces
    text = re.sub(r'[^\w\s]', '', text)  # Remove punctuation
    # Replace multiple spaces with a single space
    text = re.sub(r'\s+', ' ', text)
    text = text.strip()
    return text


# Configure logging
logging.basicConfig(level=logging.INFO,
                    format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger()

# Ensure TESSDATA_PREFIX is set
if 'TESSDATA_PREFIX' not in os.environ:
    os.environ['TESSDATA_PREFIX'] = '/usr/share/tesseract-ocr/5/tessdata/'

# Setup a global flag to indicate shutdown
shutdown = False


def preprocess_image(image_path):
    # Load the image
    image = cv2.imread(image_path, cv2.IMREAD_COLOR)

    # Convert to grayscale
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

    # Binarize the image
    _, binary = cv2.threshold(
        gray, 128, 255, cv2.THRESH_BINARY | cv2.THRESH_OTSU)

    # Denoise the image
    denoised = cv2.fastNlMeansDenoising(binary, None, 30, 7, 21)

    # Find contours
    contours, _ = cv2.findContours(
        denoised, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)

    # If no contours are found, log a message and return None
    if not contours:
        logger.warning("No text areas detected in the image.")
        return None

    # Create a mask for the detected contours
    mask = np.zeros_like(denoised)
    cv2.drawContours(mask, contours, -1, (255, 255, 255), thickness=cv2.FILLED)

    # Apply the mask to the image
    preprocessed = cv2.bitwise_and(denoised, mask)

    # Save the preprocessed image
    preprocessed_path = f"/tmp/preprocessed_{os.path.basename(image_path)}"
    cv2.imwrite(preprocessed_path, preprocessed)

    return preprocessed_path


def download_image_from_minio(client, bucket_name, object_name, file_path):
    client.fget_object(bucket_name, object_name, file_path)


def upload_image_to_minio(client, bucket_name, object_name, file_path):
    client.fput_object(bucket_name, object_name, file_path)


def update_postgres_and_elasticsearch(conn, doc_id, ocr_text, user_id):
    try:
        with conn.cursor() as cursor:
            sql = "UPDATE documents SET extract=%s WHERE id=%s"
            cursor.execute(sql, (ocr_text, doc_id))
        conn.commit()

        # Indexing in Elasticsearch
        es.index(index="documents", id=doc_id, body={
            "document_id": doc_id,
            "user_id": user_id,
            "text": ocr_text
        })

        logger.info(f"Indexed document {doc_id} in Elasticsearch")

    except Exception as e:
        logger.error(f"Failed to update PostgreSQL or Elasticsearch: {e}")


def process_message(ch, method, properties, body):
    if shutdown:
        return

    message = json.loads(body.decode())
    doc_id = message["documentID"]
    file_key = message["fileKey"]
    bucket_name = message["bucket"]
    extension = message["extension"]
    language = message['language']
    user_id = message['userID']

    file_path = f"/tmp/{os.path.basename(file_key)}.{extension}"

    logger.info(f"Received message: {message}")

    # Download image from Minio
    try:
        download_image_from_minio(
            minio_client, bucket_name, file_key, file_path)
    except Exception as e:
        logger.error(f"Failed to download image: {e}")
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
        return

    # Preprocess the image
    try:
        preprocessed_path = preprocess_image(file_path)
        if preprocessed_path is None:
            logger.warning(
                "Skipping OCR due to lack of detectable text areas.")
            ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
            return
    except Exception as e:
        logger.error(f"Failed to preprocess image: {e}")
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
        return

    # Perform OCR on the preprocessed image
    try:
        text = pytesseract.image_to_string(
            Image.open(preprocessed_path), lang=language or 'eng')
        logger.info(f"OCR Result: {text}")
    except Exception as e:
        logger.error(f"Failed to perform OCR: {e}")
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
        return
    finally:
        os.remove(file_path)  # Clean up the downloaded image
        os.remove(preprocessed_path)  # Clean up the preprocessed image

    # Update PostgreSQL with the OCR result
    if text:
        try:
            cleaned_text = clean_text(text)
            update_postgres_and_elasticsearch(
                db_conn, doc_id, cleaned_text, user_id)
        except Exception as e:
            logger.error(f"Failed to update PostgreSQL: {e}")
            ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
            return

    ch.basic_ack(delivery_tag=method.delivery_tag)
    logger.info("Done")


def signal_handler(signum, frame):
    global shutdown
    logger.info(f"Received signal {signum}, shutting down.")
    shutdown = True


if __name__ == "__main__":
    # Register signal handlers for graceful shutdown
    signal.signal(signal.SIGTERM, signal_handler)
    signal.signal(signal.SIGINT, signal_handler)

    # RabbitMQ setup
    rabbitmq_url = os.getenv("RABBITMQ_URL")
    logger.info(f"Connecting to RabbitMQ at {rabbitmq_url}")
    connection = pika.BlockingConnection(pika.URLParameters(rabbitmq_url))
    channel = connection.channel()

    channel.queue_declare(queue="extraction_queue", durable=True)
    channel.basic_qos(prefetch_count=1)
    channel.basic_consume(queue="extraction_queue",
                          on_message_callback=process_message)

    # Minio setup
    minio_client = Minio(
        os.getenv("MINIO_ENDPOINT"),
        access_key=os.getenv("MINIO_ACCESS_KEY"),
        secret_key=os.getenv("MINIO_SECRET_KEY"),
        secure=False
    )

    # PostgreSQL DB setup
    db_conn = psycopg2.connect(
        host=os.getenv("POSTGRES_HOST"),
        user=os.getenv("POSTGRES_USER"),
        password=os.getenv("POSTGRES_PASSWORD"),
        database=os.getenv("POSTGRES_DATABASE")
    )

    logger.info("Waiting for messages. To exit press CTRL+C")

    try:
        while not shutdown:
            channel.connection.process_data_events(time_limit=1)
    except Exception as e:
        logger.error(f"Error occurred: {e}")
    finally:
        channel.close()
        connection.close()
        db_conn.close()
        logger.info("RabbitMQ, PostgreSQL connection closed. Exiting.")
