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

  ingress {
    description = "HTTP for dummy nginx page"
    from_port   = 80
    to_port     = 80
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

resource "aws_eip" "box" {
  instance = aws_instance.box.id
  domain   = "vpc"

  tags = {
    Name = "bee-ci"
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
  instance_type               = "t3.micro"
  ami                         = data.aws_ami.ubuntu.id
  key_name                    = aws_key_pair.box.key_name
  vpc_security_group_ids      = [aws_security_group.box_sg.id]
  subnet_id                   = aws_subnet.public.id
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.box.name
  user_data_replace_on_change = true

  tags = {
    Name = "bee-ci"
  }

  user_data = <<-EOF
#cloud-config
package_update: true
packages:
  - curl
  - git
  - docker.io
  - nginx

write_files:
  - path: /var/www/html/index.html
    owner: www-data:www-data
    permissions: "0644"
    content: |
      <!DOCTYPE html>
      <html lang="en">
      <head>
        <meta charset="utf-8">
        <title>bee-ci box-1</title>
      </head>
      <body>
        <h1>bee-ci box-1</h1>
        <p>Dummy nginx page (PrivateLink lab).</p>
      </body>
      </html>

runcmd:
  - systemctl enable nginx
  - systemctl restart nginx
  - echo "hello from cloud-init" > /home/ubuntu/hello.txt
  - chown ubuntu:ubuntu /home/ubuntu/hello.txt
EOF
}

output "box_public_ip" {
  description = "Public IPv4 address of the EC2 box"
  value       = aws_eip.box.public_ip
}

output "box_http_url" {
  description = "Dummy nginx page on box-1"
  value       = "http://${aws_eip.box.public_ip}/"
}
