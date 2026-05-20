resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "bee-ci"
  }
}

resource "aws_lb" "main" {
  load_balancer_type = "network"
  internal           = true
  subnets            = [aws_subnet.public.id]
}

resource "aws_lb_target_group" "main" {
  port        = 80
  protocol    = "TCP"
  vpc_id      = aws_vpc.main.id
  target_type = "instance"
  preserve_client_ip = false
}

resource "aws_lb_target_group_attachment" "main" {
  target_group_arn = aws_lb_target_group.main.id
  target_id        = aws_instance.box.id
  port             = 80
}

resource "aws_lb_listener" "main" {
  load_balancer_arn = aws_lb.main.arn
  port              = 80
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.main.arn
  }
}

data "aws_caller_identity" "current" {}

resource "aws_vpc_endpoint_service" "main" {
  acceptance_required        = false
  network_load_balancer_arns = [aws_lb.main.arn]
}

resource "aws_vpc_endpoint_service_allowed_principal" "main_account" {
  vpc_endpoint_service_id = aws_vpc_endpoint_service.main.id
  principal_arn           = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
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
