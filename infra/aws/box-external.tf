resource "aws_security_group" "external_box_sg" {
  name   = "bee-ci-external-box"
  vpc_id = aws_vpc.external.id

  tags = {
    Name = "bee-ci-external"
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

resource "aws_key_pair" "external_box" {
  key_name   = "bee-ci-external-box"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILlmPPetLfPL/eTOI5wLcO3sBiY6wtjhwgm/wlQSd2LP"
}

resource "aws_instance" "external_box" {
  instance_type               = "t3.micro"
  ami                         = data.aws_ami.ubuntu.id
  key_name                    = aws_key_pair.external_box.key_name
  vpc_security_group_ids      = [aws_security_group.external_box_sg.id]
  subnet_id                   = aws_subnet.external_public.id
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.main.name
  user_data_replace_on_change = true

  tags = {
    Name = "bee-ci-external"
  }

  user_data = <<-EOF
#cloud-config
package_update: true
packages:
  - curl

runcmd:
  - echo "hello from external box" > /home/ubuntu/hello.txt
  - chown ubuntu:ubuntu /home/ubuntu/hello.txt
EOF
}

resource "aws_eip" "external_box" {
  instance = aws_instance.external_box.id
  domain   = "vpc"

  tags = {
    Name = "bee-ci-external"
  }
}

output "external_box_public_ip" {
  description = "Public IPv4 address of the external VPC EC2 box"
  value       = aws_eip.external_box.public_ip
}

output "external_box_ssh" {
  description = "SSH command for the external VPC box (use matching private key)"
  value       = "ssh -i <private-key> ubuntu@${aws_eip.external_box.public_ip}"
}
