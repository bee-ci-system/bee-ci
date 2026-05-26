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

  route {
    cidr_block                = aws_vpc.internal.cidr_block
    vpc_peering_connection_id = aws_vpc_peering_connection.internal_and_external.id
  }

  tags = {
    Name = "bee-ci-external-public"
  }
}

resource "aws_route_table_association" "external_public" {
  subnet_id      = aws_subnet.external_public.id
  route_table_id = aws_route_table.external_public.id
}

# --- PrivateLink stuff

resource "aws_security_group" "external_privatelink_endpoint" {
  name   = "bee-ci-external-privatelink-endpoint"
  vpc_id = aws_vpc.external.id

  ingress {
    description = "HTTP from external VPC"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.external.cidr_block]
  }

  egress {
    description = "Allow all outbound"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_vpc_endpoint" "main_service" {
  vpc_id             = aws_vpc.external.id
  service_name       = aws_vpc_endpoint_service.main.service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = [aws_subnet.external_public.id]
  security_group_ids = [aws_security_group.external_privatelink_endpoint.id]

}
