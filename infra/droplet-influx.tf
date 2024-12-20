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
      - 'echo "hello from cloud-init on influx! started! pwd is: $(pwd)" >> /root/hello.txt'

      # Mount the volume to /var/lib/influxdb2
      - 'mkfs.ext4 /dev/disk/by-id/scsi-0DO_Volume_influxdb-data'
      - 'mount /dev/disk/by-id/scsi-0DO_Volume_influxdb-data /var/lib/influxdb2'
      - 'echo "/dev/disk/by-id/scsi-0DO_Volume_influxdb-data /var/lib/influxdb2 ext4 defaults,nofail 0 2" | sudo tee -a /etc/fstab'

      # Install Influx
      - 'curl --silent --location -O https://repos.influxdata.com/influxdata-archive.key'
      - 'echo "943666881a1b8d9b849b74caebf02d3465d6beb716510d86a39f6c8e8dac7515 influxdata-archive.key" | sha256sum --check - && cat influxdata-archive.key | gpg --dearmor | tee /etc/apt/trusted.gpg.d/influxdata-archive.gpg > /dev/null && echo "deb [signed-by=/etc/apt/trusted.gpg.d/influxdata-archive.gpg] https://repos.influxdata.com/debian stable main" | tee /etc/apt/sources.list.d/influxdata.list'
      - 'sudo apt-get update && sudo apt-get install -y influxdb2'

      # Create the config folder
      - 'mkdir -p /etc/influxdb2'

      # Initialize InfluxDB

      - 'sudo systemctl enable influxdb'
      - 'sudo systemctl restart influxdb'
      - 'influx setup --username ${var.influxdb_user} --password ${var.influxdb_password} --token ${var.influxdb_token} --org ${var.influxdb_org} --bucket ${var.influxdb_bucket} --force'

      - 'echo "hello from cloud-init! done!" >> /root/hello.txt'

    final_message: "influx ready!"
    EOF
}

resource "digitalocean_volume" "influxdb_volume" {
  size                    = 1
  name                    = "influx-data"
  region                  = "sfo3"
  initial_filesystem_type = "ext4"
}
