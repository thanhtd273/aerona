variable "vpc_id" {
  type = string
}

variable "subnet_id" {
  description = "Subnet for Elasticache"
  type        = string
}

variable "vpc_cidr" {
  type = string
}

variable "node_type" {
  type    = string
  default = "cache.t2.micro"
}

variable "parameter_group_name" {
  type    = string
  default = "default.redis7"
}

variable "engine_version" {
  type    = string
  default = "7.1"
}
