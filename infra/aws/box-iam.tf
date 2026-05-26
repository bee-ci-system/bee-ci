resource "aws_iam_role_policy_attachment" "box_read_only" {
  role       = aws_iam_role.main.name
  policy_arn = "arn:aws:iam::aws:policy/ReadOnlyAccess"
}

resource "aws_iam_role_policy_attachment" "box_ssm_managed_instance" {
  role       = aws_iam_role.main.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role" "main" {
  name = "bee-ci-main"

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

resource "aws_iam_instance_profile" "main" {
  name = "bee-ci-main"
  role = aws_iam_role.main.name
}
