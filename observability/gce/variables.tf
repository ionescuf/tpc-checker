variable "instance_name" {
  type        = string
  description = "Name for the Google Compute instance"
}

variable "region" {
  type        = string
  description = "Region for the Google Compute instance"
}

variable "zone" {
  type        = string
  description = "Zone for the Google Compute instance"
}

variable "instance_type" {
  type        = string
  description = "Disk type of the Google Compute instance"
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
  description = "Network to attach the Google Compute instance to"
}

variable "subnet" {
  type        = string
  description = "Subnetwork to attach the Google Compute instance to"
  default     = null
}

variable "universe_domain" {
  type        = string
  description = "Universe domain for the Google Compute instance"
}

variable "image" {
  type        = string
  description = "Image family for the Google Compute instance"
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