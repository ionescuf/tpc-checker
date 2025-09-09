provider "google" {
  project         = var.project_id
  region          = var.region
  zone            = var.zone
  universe_domain = var.universe_domain
}

resource "google_container_cluster" "primary" {
  name     = var.cluster_name
  location = var.region

  network    = var.network
  subnetwork = var.subnet

  enable_autopilot = true
}

resource "google_service_account" "monitoring_sa" {
  account_id   = "monitoring-${var.cluster_name}"
  display_name = "Service Account for GKE Monitoring"
}

resource "google_project_iam_member" "monitoring_roles" {
  for_each = toset([
    "roles/monitoring.viewer",                      # Minimum for Grafana Monitoring plugin
    "roles/monitoring.metricWriter",                # If you want to write custom metrics
    "roles/logging.logWriter",                      # If logging is needed
    "roles/viewer",                                 # Optional: Full read-only access
    "roles/iam.serviceAccountOpenIdTokenCreator",   # Used by Stackdriver-Exporter for authentication 
    "roles/iam.serviceAccountTokenCreator",         # Used by Stackdriver-Exporter for authentication
    "roles/storage.admin"                           # Used by Thanos to move data to Storage bucket
  ])

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.monitoring_sa.email}"
}
