variable "name" {
  type = string
  description = "ec2 box name"
}

variable "key_name" {
    type = string
    description = "SSH public key name"
}

variable "vpc_security_group_ids" {
    type = set(string)
    description = "ids"
}

variable "subnet_id" {
  type = string
  description = "subnet id"
}

variable "instance_profile" {
  type = string
  description = "instance profile"
}
