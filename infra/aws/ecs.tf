resource "aws_ecr_repository" "main" {
  name = "bee-ci"
}

resource "aws_ecs_cluster" "main" {
  name = "bee-ci"
}

resource "aws_cloudwatch_log_group" "gh_updater" {
  name              = "/ecs/bee-ci/gh-updater"
  retention_in_days = 1
}

resource "aws_iam_role" "ecs_task_execution" {
  name = "bee-ci-ecs-task-execution"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_ecs_task_definition" "gh_updater" {
  family                   = "bee-ci-gh-updater"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512

  execution_role_arn = aws_iam_role.ecs_task_execution.arn

  container_definitions = jsonencode([
    {
      name      = "gh-updater"
      image     = "${aws_ecr_repository.main.repository_url}:gh-updater"
      essential = true

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.gh_updater.name
          awslogs-region        = "us-east-1"
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])
}

resource "aws_ecs_service" "gh_updater" {
  name            = "bee-ci-gh-updater"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.gh_updater.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [aws_subnet.internal-1.id]
    security_groups  = [aws_security_group.gh_updater.id]
    assign_public_ip = true
  }
}

resource "aws_security_group" "gh_updater" {
  name   = "bee-ci-gh-updater"
  vpc_id = aws_vpc.internal.id

  egress {
    description = "Allow outbound traffic from gh-updater"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
