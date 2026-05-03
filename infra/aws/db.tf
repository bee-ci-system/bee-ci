resource "aws_db_instance" "main" {
  engine            = "postgres"
  engine_version    = "17"
  instance_class    = "db.t4g.micro"
  allocated_storage = 8
  db_name           = "maindb"
  username          = "charlie"
  password          = "charliepass"

  # Required for terraform destroy when you do not set final_snapshot_identifier.
  skip_final_snapshot = true
}
