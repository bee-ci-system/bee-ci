resource "digitalocean_droplet" "executor" {
  name   = "vm-executor"
  region = "sfo3"
  size = "s-1vcpu-2gb" # doctl compute size list
  image  = "ubuntu-24-04-x64"

  volume_ids = []

  vpc_uuid = digitalocean_vpc.default.id

  user_data = <<-EOF
    #cloud-config
    package_update: true
    package_upgrade: true
    packages:
      - curl

     # From: https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-22-04
     runcmd:
      - sudo apt-get update
      - sudo apt-get install apt-transport-https ca-certificates curl software-properties-common
      - curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
      - echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
      - sudo apt-get update
      - apt-cache policy docker-ce

      - sudo apt install docker-ce
      - sudo systemctl status docker

      - docker pull ghcr.io/bee-ci-system/bee-ci/executor:latest

      # Postgres config
      - echo "DB_HOST=${digitalocean_database_cluster.postgres.host}" >> .executor.env
      - echo "DB_PORT=${digitalocean_database_cluster.postgres.port}" >> .executor.env
      - echo "DB_PASSWORD=${digitalocean_database_cluster.postgres.password}" >> .executor.env
      - echo "DB_NAME=${digitalocean_database_cluster.postgres.name}" >> .executor.env

      # Influx config
      - echo "INFLUXDB_URL=TODO" >> .executor.env
      - echo "INFLUXDB_ORG=${var.influxdb_org}" >> .executor.env
      - echo "INFLUXDB_BUCKET=${var.influxdb_bucket}" >> .executor.env
      - echo "INFLUXDB_TOKEN=${var.influxdb_token}" >> .executor.env

      - docker run --env-file .executor.env bee-ci-backend-executor:latest
  EOF
}
