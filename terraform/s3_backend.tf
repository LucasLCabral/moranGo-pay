terraform {
  backend "s3" {
    bucket = "morango-pay-terraform-state"
    key    = "terraform.tfstate"
    region = "sa-east-1"
  }
}
