variable "region" {
  description = "AWS region to deploy Flight Booking system"
  type        = string
  default     = "ap-southeast-1"
}

variable "postgres_password" {
  description = "Password to access PostgreSQL database"
  type        = string
}

variable "confluent_cloud_api_key" {
  type = string
}

variable "confluent_cloud_api_secret" {
  type = string
}

variable "docdb_password" {
  type = string
}

variable "public_key" {
  type      = string
  sensitive = true
}
