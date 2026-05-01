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

data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_key_pair" "box" {
  key_name   = "bee-ci-box"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILlmPPetLfPL/eTOI5wLcO3sBiY6wtjhwgm/wlQSd2LP"
}

resource "aws_instance" "box" {
  availability_zone = "us-east-1a"
  instance_type     = "t3.micro"
  ami               = data.aws_ami.ubuntu.id
  key_name          = aws_key_pair.box.key_name
}

output "box_public_ip" {
  description = "Public IPv4 address of the EC2 box"
  value       = aws_instance.box.public_ip
}
