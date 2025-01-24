#!/bin/sh

# Run migrations
echo "Running migrations..."
go run cmd/migrate/main.go up

# Start the application
echo "Starting the application..."
exec ./bin/suraksheet
