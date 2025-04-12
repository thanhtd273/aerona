
variable "db_username" {
  type    = string
  default = "admin"
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "subnet_id" {
  type = string

}

variable "vpc_id" {
  type = string
}

variable "vpc_cidr" {
  type = string
}

variable "db_instance_class" {
  type    = string
  default = "db.t3.medium"
}
