

This directory provides everything you need to deploy a robust observability stack on Google Kubernetes Engine (GKE) using Helm and Terraform.

### Contents

- **Helm values files** for [kube-prometheus-stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack) and [prometheus-stackdriver-exporter](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus-stackdriver-exporter)
- **Terraform configuration** (in the `terraform/` folder) for provisioning GKE and related resources

### 1. Infrastructure Provisioning with Terraform

The `terraform/` folder contains Terraform code to provision your GKE cluster and supporting Google Cloud resources.

**Typical resources managed:**
- GKE cluster and node pools
- Service accounts and IAM roles
- Networking (VPC, subnets, firewall rules)
- Optionally, Google Cloud Monitoring and logging integrations

**How to use:**
1. Edit `terraform/variables.tf` and/or `terraform.tfvars` to match your project and region.
2. Initialize and apply:
   ```sh
   cd terraform
   terraform init
   terraform plan -var-file="variables.tfvars"
   terraform apply -var-file="variables.tfvars"
   ```
3. After completion, configure `kubectl` to use the new cluster (see Terraform output or use `gcloud container clusters get-credentials ...`).

## 2. Deploy kube-prometheus-stack

This chart deploys Prometheus, Grafana, Alertmanager, and related monitoring components.

**Configuration:**  
The file [`values-prometheus-stack.yaml`](./values-prometheus-stack.yaml) customizes the deployment:

- **Prometheus, Grafana, and Alertmanager** are exposed as LoadBalancer services.
- **ServiceMonitor** for the Stackdriver exporter is included, so Prometheus scrapes metrics from it.

**Install:**
```sh
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/kube-prometheus-stack --namespace monitoring --create-namespace -f values-prometheus-stack.yaml
```

### 3. Deploy prometheus-stackdriver-exporter

This Helm chart exposes Google Cloud Monitoring (Stackdriver) metrics to Prometheus.

**Install:**
```sh
helm install stackdriver-exporter prometheus-community/prometheus-stackdriver-exporter --namespace monitoring --set stackdriver.projectId="eu0:project-name" -f values-stackdriver-exporter.yaml
```

### 4. Accessing the Services

After deployment, get the external IPs for the services:

```sh
kubectl get svc -n monitoring
```

- **Prometheus:** `http://<EXTERNAL-IP>:9090`
- **Grafana:** `http://<EXTERNAL-IP>:3000` (default login: `admin` / `prom-operator`)
- **Alertmanager:** `http://<EXTERNAL-IP>:9093`

### 5. Customization

- Edit `values-prometheus-stack.yaml` to adjust service types, ports, ServiceMonitors, and enabled components.
- For advanced Stackdriver exporter configuration adjust `values-stackdriver-exporter.yaml`.

For a comprehensive list of configuration options, please refer to the official `values.yaml` files of both Helm charts, which include example configurations:
- [Prometheus Stackdriver Exporter `values.yaml`](https://github.com/prometheus-community/helm-charts/blob/main/charts/prometheus-stackdriver-exporter/values.yaml)
- [Kube Prometheus Stack `values.yaml`](https://github.com/prometheus-community/helm-charts/blob/main/charts/kube-prometheus-stack/values.yaml)

### 6. Connecting Stackdriver Exporter to the Monitoring API

To enable Stackdriver Exporter to connect to the TPC Monitoring API, you must link the Kubernetes service account to the TPC service account.
To link the accounts, set the annotation `iam.gke.io/gcp-service-account` on the Kubernetes service account. The value of this annotation should be the name of the TPC service account.

We also need to bind the roles/iam.workloadIdentityUser role to the Kubernetes service account.

Example:
```sh
gcloud iam service-accounts add-iam-policy-binding \
  "monitoring-observability-gke@observability-epam.eu0.iam.gserviceaccount.com" \
  --role="roles/iam.workloadIdentityUser" \
  --member="serviceAccount:observability-epam.eu0.svc.id.goog[monitoring/stackdriver-exporter]" \
  --project="eu0:observability-epam"
```
- gcloud iam service-accounts add-iam-policy-binding: This is the core command that modifies the IAM policy of a service account.

- "monitoring-observability-gke@observability-epam.eu0.iam.gserviceaccount.com": This is the Google Cloud service account that you're granting permissions to.

- "roles/iam.workloadIdentityUser": This specifies the role being granted. The Workload Identity User role allows the Kubernetes service account to impersonate the Google Cloud service account.

- member="serviceAccount:observability-epam.eu0.svc.id.goog[monitoring/stackdriver-exporter]": This is the member receiving the role. It specifies the Kubernetes service account (stackdriver-exporter) in a specific namespace (monitoring) within your Google Cloud project.

- project="eu0:observability-epam": This is the Google Cloud project ID where the Google Cloud service account is located.

This allows the Kubernetes service account to impersonate the TPC service account, granting it the necesarry permissions.

### 7. Uninstall

To remove the deployments:

```sh
helm uninstall prometheus-stack -n monitoring
helm uninstall stackdriver-exporter -n monitoring
```

### References

- [kube-prometheus-stack Helm Chart](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)
- [prometheus-stackdriver-exporter Helm Chart](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus-stackdriver-exporter)
- [Terraform](https://www.terraform.io/)

