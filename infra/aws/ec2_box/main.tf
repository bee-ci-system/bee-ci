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

resource "aws_instance" "small_ubuntu_box" {
  instance_type               = "t3.micro"
  ami                         = data.aws_ami.ubuntu.id
  key_name                    = var.key_name
  vpc_security_group_ids      = var.vpc_security_group_ids
  subnet_id                   = var.subnet_id
  associate_public_ip_address = true
  iam_instance_profile        = var.instance_profile

  tags = {
    Name = var.name
  }

  user_data = <<EOF
    #cloud-config
    package_update: true
    packages:
      - curl
      - git
      - docker.io

    runcmd:
      - echo "hello from cloud-init on box ${var.name}" > /home/ubuntu/hello.txt
      - chown ubuntu:ubuntu /home/ubuntu/hello.txt
  EOF
}