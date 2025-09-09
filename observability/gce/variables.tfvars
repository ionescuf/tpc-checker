# TPC project ID
project_id      = "eu0:observability-epam"

# Name for the monitoring VM instance
instance_name   = "observability"

# Machine type for the VM
instance_type   = "c3-standard-4"

# TPC region and zone
region          = "u-germany-northeast1"
zone            = "u-germany-northeast1-a"

# VPC network and subnet (full resource path recommended)
network         = "projects/eu0:observability-epam/global/networks/observability-epam-vpc"
subnet          = "projects/eu0:observability-epam/regions/u-germany-northeast1/subnetworks/observability-epam-vpc"

# Persistent disk size in GB
disk_size_gb    = 10

# Universe domain for TPC APIs
universe_domain = "apis-berlin-build0.goog"

# VM image to use (full resource path)
image           = "projects/eu0-system:debian-cloud/global/images/debian-12--tpc-20250611-2318"

# List of allowed source IP ranges for firewall
source_ranges   = ["195.56.119.208/28"]

# List of allowed ports for firewall rules
allowed_ports   = ["80", "443", "3000", "9090", "9093", "9100"]