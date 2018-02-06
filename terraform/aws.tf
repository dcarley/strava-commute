terraform {
  backend "s3" {
    key = "terraform.tfstate"
  }
}

provider "aws" {}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}
