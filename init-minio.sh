#!/bin/sh

# Start Minio server in the background
minio server /data --console-address ":9001" &

# Capture PID of Minio server
MINIO_PID=$!

# Wait until Minio is up and running
until curl -s http://minio:9000/minio/health/live
do
  echo "Waiting for Minio..."
  sleep 5
done

# Configure mc with the Minio credentials
mc alias set myminio http://minio:9000 ${MINIO_ROOT_USER} ${MINIO_ROOT_PASSWORD}

# Create the bucket
echo "Creating bucket 'suraksheet'..."
mc mb myminio/suraksheet

echo "Bucket 'suraksheet' created."

# Wait for termination signal
trap "echo 'Stopping Minio...'; kill $MINIO_PID; wait $MINIO_PID; exit 0" SIGTERM

# Keep the script running to avoid container exit
wait
