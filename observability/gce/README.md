The observability folder contains scripts and Terraform code that allow to launch a new VM in TPC/GCP, then deploy Prometheus and Graphaha there.

The intended audience is users of TPC who may need an example of such setup.

Not intended for production use, this is for demonstration purposes only.

# Observability GCE Stack Deployment Guide

This guide explains how to deploy the monitoring infrastructure using Terraform and Docker Compose in the `tpc-checker/observability/gce` folder.

## Folder Structure

```
tpc-checker/observability/gce/
├── alertmanager/
│   └── config.yml
├── app/
│   └── Dockerfile
│   └── ...
├── grafana/
│   └── provisioning/
│       └── datasources/
│           └── datasource.yml
│       └── dashboards/
├── prometheus/
│   └── prometheus.yml
├── stackdriver-exporter/
│   └── Dockerfile
├── docker-compose.yaml
├── main.tf
├── outputs.tf
├── variables.tf
├── variables.tfvars
├── startup-script.sh
```

### Deploying Infrastructure with Terraform

### Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) installed
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) installed and authenticated (`gcloud auth application-default login`)
- Go (version >1.23.0) installed

### Steps

1. Configure Variables

   Edit `variables.tfvars` 

2. Initialize Terraform**

   ```sh
   terraform init
   ```

3. Review the Plan**

   ```sh
   terraform plan -var-file="variables.tfvars"
   ```

4. Apply the Configuration**
   ```sh
   terraform apply -var-file="variables.tfvars"
   ```

   This will create:
   - Create a Google Compute Engine (GCE) VM instance.
   - Configure the VM with a startup script (startup-script.sh).
   - Set up firewall rules to allow access to Prometheus, Grafana, Alertmanager, and Node Exporter ports (e.g., 3000, 9090, 9093, 9100).
   - Create a service account and assign necessary IAM roles for monitoring.


### Deploying the Applications with Docker Compose

1. SSH into the VM

   ```sh
   gcloud compute ssh <your-vm-name> --zone <your-zone>
   ```

2. Clone the repository to the VM

   ```sh
   git clone https://github.com/epam/tpc-checker.git
   ```

3. Clone the stackdriver_exporter repo to the VM

   ```sh
   git clone https://github.com/prometheus-community/stackdriver_exporter.git
   ```

4. Build the stackdriver-exporter image

   ```sh
   cd stackdriver_exporter
   make build
   sudo docker build -t tpc-stackdriver-exporter .
   ```

6. Start the containers

   ```sh
   cd ~/tpc-checker/observability/gce
   sudo docker-compose up -d --build
   ```

### Accessing the Services

- **Grafana:** http://<VM_EXTERNAL_IP>:3000  
  Default credentials: `admin` / `admin`
- **Prometheus:** http://<VM_EXTERNAL_IP>:9090
- **Alertmanager:** http://<VM_EXTERNAL_IP>:9093
- **Node Exporter:** http://<VM_EXTERNAL_IP>:9100


### Configuration files

- Prometheus config: `prometheus/prometheus.yml`
- Grafana provisioning: 
    - `grafana/provisioning/datasources/datasource.yml`
    - `grafana/provisioning/dashboards/dashboards.yml`
- Alertmanager config: `alertmanager/config.yml`

You can customize these files before deployment.


### Troubleshooting 

Use `docker ps` and `docker logs <container>` on the VM to debug container issues.

### Troubleshooting stackdriver-exporter Build Issues
You might encounter the following error when building stackdriver-exporter:

```sh
curl -s -L https://github.com/prometheus/promu/releases/download/v0.17.0/promu-0.17.0.-.tar.gz | tar -xvzf - -C /tmp/tmp.2JdQWjeCZP
gzip: stdin: not in gzip format
tar: Child returned status 1
tar: Error is not recoverable: exiting now
make: *** [Makefile.common:246: /bin/promu] Error 2
Cause: This error indicates that the operating system type is not recognized in the promu download URL. The dash (-) in the filename needs to be replaced with your actual OS and architecture.
```
Solution: Replace the `-` in the URL with your system's correct OS and architecture (e.g., `linux-amd64`, `darwin-amd64`) and retry `make build`. 
For example, for Linux AMD64:

```sh
curl -s -L https://github.com/prometheus/promu/releases/download/v0.17.0/promu-0.17.0.linux-amd64.tar.gz | tar -xvzf - -C /tmp/tmp.2JdQWjeCZP
```

### IAM Service Account Key
The `stackdriver-exporter` requires a key pair for an IAM Service Account to authenticate with Google Cloud.

Steps to Create and Download the Key:
1. Navigate to GCP IAM & Admin > Service Accounts.
2. Select your desired Service Account (Grafana Monitoring Service Account).
3. Go to the Keys tab.
4. Click Add Key > Create new key.
5. Choose JSON as the key type and click Create.
6. Download the generated JSON file and save it securely on your VM where the docker-compose file is referencing it at line 92.


