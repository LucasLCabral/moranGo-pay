terraform {
  required_version = ">= 1.0"
  
  backend "s3" {
    bucket = "morango-pay-terraform-state"
    key    = "terraform.tfstate"
    region = "sa-east-1"
  }
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region = "sa-east-1"
}