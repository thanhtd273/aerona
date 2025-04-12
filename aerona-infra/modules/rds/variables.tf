variable "vpc_cidr" {
  description = "CIDR block for Postgres"
  type        = string
}

variable "vpc_id" {
  description = "VPC for PostgreSQL"
  type        = string
}

variable "subnet_id" {
  description = "VPC Subnet ID"
  type        = string
}

variable "db_engine_version" {
  description = "Engine version for Postgres"
  type        = string
  default     = "16.8-R1"
}

variable "db_instance_class" {
  description = "DB instance for RDS"
  type        = string
  default     = "db.t4g.micro"
}

variable "db_storage" {
  description = "Storage size (GB)"
  type        = number
  default     = 10
}

variable "db_username" {
  description = "Admin username"
  type        = string
  default     = "admin"
}

variable "db_password" {
  description = "Admin password"
  type        = string
  sensitive   = true
}
