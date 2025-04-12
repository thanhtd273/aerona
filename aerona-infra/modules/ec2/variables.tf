variable "vpc_id" {
  description = "ID of the VPC"
  type        = string
}

variable "subnet_id" {
  description = "EC2 subnet"
  type        = string
}

variable "ami" {
  description = "Software image for EC2 instance"
  type        = string
  default     = "ami-065a492fef70f84b1"
}

variable "public_key" {
  type = string
}
