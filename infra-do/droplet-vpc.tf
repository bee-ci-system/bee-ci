resource "digitalocean_vpc" "default" {
  name     = "vpc-main"
  region   = "sfo3"
  ip_range = "10.0.0.0/16"
}

resource "digitalocean_ssh_key" "default" {
  name       = "main key"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILlmPPetLfPL/eTOI5wLcO3sBiY6wtjhwgm/wlQSd2LP"
  # "bee-ci droplets key" in 1Password
}
