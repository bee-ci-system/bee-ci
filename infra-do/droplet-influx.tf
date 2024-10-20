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

# Ubuntu and Debian
  # Add the InfluxData key to verify downloads and add the repository




    # Install influxdb



    # From: https://docs.influxdata.com/influxdb/v2/install
    runcmd:
      - echo "hello from cloud-init! started! pwd is: $(pwd)" >> /root/hello.txt

      # Install InfluxDB
      - curl --silent --location -O https://repos.influxdata.com/influxdata-archive.key
      - echo "943666881a1b8d9b849b74caebf02d3465d6beb716510d86a39f6c8e8dac7515 influxdata-archive.key" | sha256sum --check - && cat influxdata-archive.key | gpg --dearmor | tee /etc/apt/trusted.gpg.d/influxdata-archive.gpg > /dev/null && echo 'deb [signed-by=/etc/apt/trusted.gpg.d/influxdata-archive.gpg] https://repos.influxdata.com/debian stable main' | tee /etc/apt/sources.list.d/influxdata.list
      - sudo apt-get update && sudo apt-get install -y influxdb2

      # Mount the volume to /var/lib/influxdb2
      - mkfs.ext4 /dev/disk/by-id/scsi-0DO_Volume_influxdb-data
      - mount /dev/disk/by-id/scsi-0DO_Volume_influxdb-data /var/lib/influxdb2
      - echo "/dev/disk/by-id/scsi-0DO_Volume_influxdb-data /var/lib/influxdb2 ext4 defaults,nofail 0 2" | sudo tee -a /etc/fstab

      # Create the config folder
      - mkdir -p /etc/influxdb2

      # Initialize InfluxDB with the provided environment variables
      - INFLUXDB_INIT_USERNAME=beeci
      - INFLUXDB_INIT_PASSWORD=${var.influxdb_password}
      - INFLUXDB_INIT_ADMIN_TOKEN=${var.influxdb_token}
      - INFLUXDB_INIT_ORG=${var.influxdb_org}
      - INFLUXDB_INIT_BUCKET=${var.influxdb_bucket}

      - influx setup --username beeci --password ${var.influxdb_password} --token ${var.influxdb_token} --org ${var.influxdb_org} --bucket ${var.influxdb_bucket} --force

      # Restart InfluxDB to apply storage and initialization changes
      - sudo systemctl enable influxdb
      - sudo systemctl restart influxdb
      - echo "hello from cloud-init! done!" >> /root/hello.txt
    EOF
}

resource "digitalocean_volume" "influxdb_volume" {
  size                    = 1
  name                    = "influx-data"
  region                  = "sfo3"
  initial_filesystem_type = "ext4"
}

# To check if cloud-init completed successfully, see:
# https://www.digitalocean.com/community/questions/how-to-make-sure-that-cloud-init-finished-running
