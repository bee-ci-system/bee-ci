resource "aws_vpc" "external" {
  cidr_block           = "10.1.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "bee-ci-external"
  }
}

resource "aws_subnet" "external_public" {
  vpc_id                  = aws_vpc.external.id
  cidr_block              = "10.1.1.0/24"
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = true

  tags = {
    Name = "bee-ci-external-public"
  }
}

resource "aws_internet_gateway" "external" {
  vpc_id = aws_vpc.external.id

  tags = {
    Name = "bee-ci-external"
  }
}

resource "aws_route_table" "external_public" {
  vpc_id = aws_vpc.external.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.external.id
  }

  tags = {
    Name = "bee-ci-external-public"
  }
}

resource "aws_route_table_association" "external_public" {
  subnet_id      = aws_subnet.external_public.id
  route_table_id = aws_route_table.external_public.id
}
