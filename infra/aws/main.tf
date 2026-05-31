terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_vpc_peering_connection" "internal_and_external" {
  vpc_id      = aws_vpc.internal.id
  peer_vpc_id = aws_vpc.external.id
  auto_accept = true

  tags = {
    Name = "bee-ci-internal-external"
  }
}
