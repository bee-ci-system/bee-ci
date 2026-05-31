locals {
  instances = [aws_instance.box_internal.id, aws_instance.external_box.id]
}

output "instances" {
  description = "List of instances"
  value       = local.instances
}

resource "local_file" "my_json_file" {
  filename = "${path.module}/instances.json"
  content  = jsonencode(local.instances)
}
