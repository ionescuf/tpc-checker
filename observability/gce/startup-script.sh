#!/bin/bash

# Install Docker and Docker Compose
apt-get update
apt-get install -y docker.io docker-compose make
systemctl enable docker
systemctl start docker

# Optional: Add default user to docker group
usermod -aG docker $USER

# Mount persistent disk
sudo mkfs.ext4 /dev/nvme0n2
mkdir -p /mnt/docker-data/{grafana,prometheus} 
sudo mount /dev/nvme0n2 /mnt/docker-data
cd /mnt/docker-data

# # Create Docker volumes for Grafana and Prometheus
docker volume create --opt type=none --opt device=/mnt/docker-data/grafana --opt o=bind grafana_data
docker volume create --opt type=none --opt device=/mnt/docker-data/prometheus --opt o=bind prometheus_data
