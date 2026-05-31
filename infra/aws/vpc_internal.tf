resource "aws_vpc" "internal" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "bee-ci"
  }
}

resource "aws_lb" "internal" {
  load_balancer_type = "network"
  internal           = true
  subnets            = [aws_subnet.internal-1.id]
}

resource "aws_lb_target_group" "internal" {
  port               = 80
  protocol           = "TCP"
  vpc_id             = aws_vpc.internal.id
  target_type        = "instance"
  preserve_client_ip = false
}

resource "aws_lb_target_group_attachment" "internal" {
  target_group_arn = aws_lb_target_group.internal.id
  target_id        = aws_instance.box_internal.id
  port             = 80
}

resource "aws_lb_listener" "internal" {
  load_balancer_arn = aws_lb.internal.arn
  port              = 80
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.internal.arn
  }
}

data "aws_caller_identity" "current" {}

resource "aws_vpc_endpoint_service" "main" {
  acceptance_required        = false
  network_load_balancer_arns = [aws_lb.internal.arn]
}

resource "aws_vpc_endpoint_service_allowed_principal" "main_account" {
  vpc_endpoint_service_id = aws_vpc_endpoint_service.main.id
  principal_arn           = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
}

resource "aws_subnet" "internal-1" {
  vpc_id                  = aws_vpc.internal.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = true

  tags = {
    Name = "bee-ci-internal-1"
  }
}

resource "aws_internet_gateway" "internal" {
  vpc_id = aws_vpc.internal.id

  tags = {
    Name = "bee-ci"
  }
}

resource "aws_route_table" "internal" {
  vpc_id = aws_vpc.internal.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.internal.id
  }

  route {
    cidr_block                = aws_vpc.external.cidr_block
    vpc_peering_connection_id = aws_vpc_peering_connection.internal_and_external.id
  }

  tags = {
    Name = "bee-ci-internal"
  }
}

resource "aws_route_table_association" "internal" {
  subnet_id      = aws_subnet.internal-1.id
  route_table_id = aws_route_table.internal.id
}
