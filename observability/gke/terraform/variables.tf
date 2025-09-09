variable "cluster_name" {
  type        = string
  description = "Name for the Google GKE Cluster"
}

variable "region" {
  type        = string
  description = "Region for the Google GKE Cluster"
}

variable "zone" {
  type        = string
  description = "Zone for the Google GKE Cluster"
}

variable "instance_type" {
  type        = string
  description = "Disk type of the Google GKE Cluster"
}

variable "project_id" {
  type        = string
  description = "Google Cloud project ID"
}

variable "disk_size_gb" {
  type        = number
  description = "Size of the persistent disk in GB"
  default     = 30
}

variable "network" {
  type        = string
  description = "Network to attach the Google GKE Cluster to"
}

variable "subnet" {
  type        = string
  description = "Subnetwork to attach the Google GKE Cluster to"
  default     = null
}

variable "universe_domain" {
  type        = string
  description = "Universe domain for the Google GKE Cluster"
}

variable "image" {
  type        = string
  description = "Image family for the Google GKE Cluster"
}

variable "node_pool" {
  type        = string
  description = "Name of the GKE Cluster node pool"
}

variable "source_ranges" {
  type        = list(string)
  description = "List of source IP ranges allowed by the firewall"
}

variable "allowed_ports" {
  type        = list(string)
  description = "List of allowed ports for the firewall rule"
  default     = ["80", "443", "3000", "9090", "9093", "9100"]
}
