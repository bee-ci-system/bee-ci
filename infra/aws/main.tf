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

resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "bee-ci"
  }
}

resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = true

  tags = {
    Name = "bee-ci-public"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "bee-ci"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "bee-ci-public"
  }
}

resource "aws_route_table_association" "public" {
  subnet_id      = aws_subnet.public.id
  route_table_id = aws_route_table.public.id
}

resource "aws_security_group" "box_sg" {
  name   = "bee-ci-box"
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "bee-ci"
  }

  ingress {
    description = "SSH from anywhere"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "Allow all outbound"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

locals {
  names = toset(["1", "2", "3"])
}

resource "aws_eip" "box_eip" {
  for_each = local.names

  instance = aws_instance.box[each.value].id
  domain   = "vpc"

  tags = {
    Name = "bee-ci-${each.value}"
  }
}

resource "aws_key_pair" "box" {
  key_name   = "bee-ci-box"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILlmPPetLfPL/eTOI5wLcO3sBiY6wtjhwgm/wlQSd2LP"
}

resource "aws_iam_role" "box" {
  name = "bee-ci-box"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })

  tags = {
    Name = "bee-ci"
  }
}

resource "aws_iam_role_policy_attachment" "box_read_only" {
  role       = aws_iam_role.box.name
  policy_arn = "arn:aws:iam::aws:policy/ReadOnlyAccess"
}

resource "aws_iam_instance_profile" "box" {
  name = "bee-ci-box"
  role = aws_iam_role.box.name
}

resource "aws_instance" "box" {
  for_each = local.names

  instance_type               = "t3.micro"
  ami                         = data.aws_ami.ubuntu.id
  key_name                    = aws_key_pair.box.key_name
  vpc_security_group_ids      = [aws_security_group.box_sg.id]
  subnet_id                   = aws_subnet.public.id
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.box.name

  tags = {
    Name = "bee-ci-${each.value}"
  }

  user_data = <<EOF
    #cloud-config
    package_update: true
    packages:
      - curl
      - git
      - docker.io

    runcmd:
      - echo "hello from cloud-init" > /home/ubuntu/hello.txt
      - chown ubuntu:ubuntu /home/ubuntu/hello.txt
  EOF
}

output "box_public_ip" {
  description = "Public IPv4 address of the EC2 box"
  value = {
    for key, eip in aws_eip.box_eip : "IP of box ${key}" => eip.public_ip
  }
}
