resource "digitalocean_droplet" "executor" {
  name          = "vm-executor"
  region        = "sfo3"
  image         = "ubuntu-24-04-x64"
  size = "s-1vcpu-2gb" # doctl compute size list
  volume_ids = []
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

    # From: https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-22-04
    runcmd:
      - 'echo "hello from cloud-init on executor! started! pwd is: $(pwd)" >> /root/hello.txt'
      - 'sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common'
      - 'curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg'
      - 'echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null'
      - 'sudo apt-get update'
      - 'apt-cache policy docker-ce'

      - 'sudo apt-get install -y docker-ce'
      - 'sudo systemctl status docker'

      - 'docker pull ghcr.io/bee-ci-system/bee-ci/executor:latest'

      # Postgres config
      - 'echo "DB_HOST=${digitalocean_database_cluster.postgres.host}" >> /root/.executor.env'
      - 'echo "DB_PORT=${digitalocean_database_cluster.postgres.port}" >> /root/.executor.env'
      - 'echo "DB_PASSWORD=${digitalocean_database_cluster.postgres.password}" >> /root/.executor.env'
      - 'echo "DB_NAME=${digitalocean_database_cluster.postgres.name}" >> /root/.executor.env'

      # Influx config
      - 'echo "INFLUXDB_URL=http://${digitalocean_droplet.influxdb.ipv4_address_private}:8086" >> /root/.executor.env'
      - 'echo "INFLUXDB_ORG=${var.influxdb_org}" >> /root/.executor.env'
      - 'echo "INFLUXDB_BUCKET=${var.influxdb_bucket}" >> /root/.executor.env'
      - 'echo "INFLUXDB_TOKEN=${var.influxdb_token}" >> /root/.executor.env'

      - 'docker run --env-file /root/.executor.env bee-ci-backend-executor:latest'
      - 'echo "hello from cloud-init! done!" >> /root/hello.txt'

    final_message: "executor ready!"
    EOF
}
