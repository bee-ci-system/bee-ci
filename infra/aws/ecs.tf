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
