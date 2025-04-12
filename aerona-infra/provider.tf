terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }

    confluent = {
      source  = "confluentinc/confluent"
      version = "2.24.0"
    }

  }
  required_version = ">= 1.2.0"
}

provider "aws" {
  region = var.region
}

provider "confluent" {
  cloud_api_key    = var.confluent_cloud_api_key
  cloud_api_secret = var.confluent_cloud_api_secret
}
