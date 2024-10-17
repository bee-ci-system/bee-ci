resource "digitalocean_vpc" "default" {
  name     = "vpc-main"
  region   = "sfo3"
  ip_range = "10.0.0.0/16"
}
