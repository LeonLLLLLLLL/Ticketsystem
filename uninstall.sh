#!/bin/bash

echo "Stopping and removing the project's containers..."
CURRENT_DIR_NAME=$(basename "$PWD")
docker-compose down --volumes

# Get image names from docker-compose and remove them
echo "Removing project's images..."
docker rmi "${CURRENT_DIR_NAME}_frontend"
docker rmi "${CURRENT_DIR_NAME}_address_module_backend"
#docker rmi -f $(docker images -q $(docker-compose config | awk '/image:/ {print $2}'))

echo "Removing project's volumes..."
docker volume rm -f mysql_data 2>/dev/null

echo "Removing project's network..."
docker network rm mynetwork 2>/dev/null

echo "Project-specific Docker cleanup complete."
