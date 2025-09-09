provider "google" {
  project         = var.project_id
  region          = var.region
  zone            = var.zone
  universe_domain = var.universe_domain
}

resource "google_compute_firewall" "default" {
  name    = "allow-http-https"
  network = var.network

  allow {
    protocol = "tcp"
    ports    = var.allowed_ports
  }

  source_ranges = var.source_ranges
  target_tags   = ["web-server", "monitoring"]
}

resource "google_compute_disk" "persistent_disk" {
  name = "${var.instance_name}-disk"
  type = "hyperdisk-balanced"
  zone = var.zone
  size = var.disk_size_gb
}

resource "google_compute_instance" "monitoring_vm" {
  name         = var.instance_name
  machine_type = var.instance_type
  zone         = var.zone
  tags         = ["web-server"]

  boot_disk {
    initialize_params {
      image = var.image
    }
  }

  attached_disk {
    source      = google_compute_disk.persistent_disk.id
    device_name = "data-disk"
  }

  network_interface {
    network    = var.network
    subnetwork = var.subnet
    access_config {}
  }

  metadata_startup_script = file("startup-script.sh")

  service_account {
    email  = google_service_account.grafana_sa.email
    scopes = ["https://www.googleapis.com/auth/cloud-platform"]
  }
}

resource "google_service_account" "grafana_sa" {
  account_id   = "grafana-monitoring-sa"
  display_name = "Grafana Monitoring Service Account"
}

resource "google_project_iam_member" "monitoring_roles" {
  for_each = toset([
    "roles/monitoring.viewer",       # Minimum for Grafana Monitoring plugin
    "roles/monitoring.metricWriter", # If you want to write custom metrics
    "roles/logging.logWriter",       # If logging is needed
    "roles/viewer"                   # Optional: Full read-only access
  ])
  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.grafana_sa.email}"
}
