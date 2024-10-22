resource "digitalocean_droplet" "influxdb" {
  name          = "vm-influx"
  region        = "sfo3"
  image         = "ubuntu-24-04-x64"
  size = "s-1vcpu-512mb-10gb" # doctl compute size list
  volume_ids = [digitalocean_volume.influxdb_volume.id]
  tags = []
  vpc_uuid      = digitalocean_vpc.default.id
  ssh_keys = [digitalocean_ssh_key.default.fingerprint]
  droplet_agent = true

  user_data = <<-EOF
    #cloud-config
    package_update: true
    package_upgrade: true
    packages:
      - curl
      - micro
      - bat

    # From: https://docs.influxdata.com/influxdb/v2/install
    runcmd:
      - echo "hello from cloud-init on influx! started! pwd is: $(pwd)" >> /root/hello.txt

    final_message: "influx ready!"
    EOF
}

resource "digitalocean_volume" "influxdb_volume" {
  size                    = 1
  name                    = "influx-data"
  region                  = "sfo3"
  initial_filesystem_type = "ext4"
}
