terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "2.41.0"
    }
  }

  cloud {
    organization = "bpacia"

    workspaces {
      project = "bee-ci"
      name    = "bee-ci-workspace"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_project" "project" {
  name        = "bee-ci-proj-tf"
  description = "A simple container-based CI system"
  environment = "Development"
  resources = [
    digitalocean_app.app.urn,
    digitalocean_database_cluster.postgres.urn,
    digitalocean_database_cluster.redis.urn,
    digitalocean_domain.default.urn,
    # digitalocean_droplet.influxdb.urn,
    digitalocean_volume.influxdb_volume.urn,
    # digitalocean_domain.main.urn,
  ]
}

resource "digitalocean_app" "app" {

  spec {
    name   = "bee-ci-tf"
    region = "sfo"

    # domain {
    #   # name = "beeci-backend.ondigitalocean.app"
    #   type = "DEFAULT"
    # }

    domain {
      name = "beeci-backend.karolak.cc"
      type = "PRIMARY"
    }

    ingress {

      /*
      rule {
        component {
          name = "frontend"
        }
        match {
          path {
            prefix = "/"
          }
        }
      }
       */

      rule {
        component {
          name = "server"
        }
        match {
          path {
            prefix = "/backend"
          }
        }
      }
    }

    /*
    database {
      engine       = "PG"
      name         = digitalocean_database_db.postgres.name
      db_name      = digitalocean_database_db.postgres.name
      cluster_name = digitalocean_database_cluster.postgres.name
      production   = true
    }

    database {
      engine       = "REDIS"
      name         = digitalocean_database_db.redis.name
      db_name      = digitalocean_database_db.postgres.name
      cluster_name = digitalocean_database_cluster.postgres.name
      production   = true
    }
     */

    /*
    service {
      name               = "frontend"
      instance_count     = 1
      instance_size_slug = "apps-s-1vcpu-0.5gb" # doctl apps tier instance-size list

      http_port = 3000

      github {
        repo           = "bee-ci-system/bee-ci"
        branch         = "master"
        deploy_on_push = true
      }

      source_dir      = "./frontend"
      dockerfile_path = "./frontend/Dockerfile"

      health_check {
        http_path             = "/"
        initial_delay_seconds = 10
        period_seconds        = 5
        timeout_seconds       = 1
        success_threshold     = 3
        failure_threshold     = 3
      }
    }
    */

    service {
      name               = "server"
      instance_count     = 3
      instance_size_slug = "apps-s-1vcpu-1gb" # doctl apps tier instance-size list

      dynamic "env" {
        for_each = local.env_vars
        content {
          key   = env.value.key
          value = env.value.value
          scope = env.value.scope
          type = lookup(env.value, "type", null)
        }
      }

      http_port = 8080

      github {
        repo           = "bee-ci-system/bee-ci"
        branch         = "master"
        deploy_on_push = true
      }

      source_dir      = "./backend"
      dockerfile_path = "./backend/server.dockerfile"

      health_check {
        http_path             = "/"
        initial_delay_seconds = 10
        period_seconds        = 5
        timeout_seconds       = 1
        success_threshold     = 3
        failure_threshold     = 3
      }
    }


    worker {
      name               = "gh-updater"
      instance_count     = 1
      instance_size_slug = "apps-s-1vcpu-0.5gb" # doctl apps tier instance-size list

      dynamic "env" {
        for_each = local.env_vars
        content {
          key   = env.value.key
          value = env.value.value
          scope = env.value.scope
          type = lookup(env.value, "type", null)
        }
      }

      github {
        repo           = "bee-ci-system/bee-ci"
        branch         = "master"
        deploy_on_push = true
      }

      source_dir      = "./backend"
      dockerfile_path = "./backend/gh-updater.dockerfile"
    }
  }
}


resource "digitalocean_container_registry" "default" {
  name   = "bee-ci-container-registry"
  subscription_tier_slug = "basic" # $5/month
  region = "sfo2"
}

resource "digitalocean_container_registry_docker_credentials" "default" {
  registry_name = "bee-ci-container-registry"
}


resource "digitalocean_volume" "influxdb_volume" {
  size = 1 # GB
  name                    = "influxdb-data"
  region                  = "sfo3"
  initial_filesystem_type = "ext4"
}

/*
resource "digitalocean_droplet" "influxdb" {
  name   = "influxdb-server"
  region = "sfo3"
  size   = "s-1vcpu-512mb-10gb" # doctl compute size list
  image  = "ubuntu-24-04-x64"

  volume_ids = [digitalocean_volume.influxdb_volume.id]

  user_data = <<-EOF
    #cloud-config
    package_update: true
    package_upgrade: true
    packages:
      - curl

     runcmd:
      - curl -sL https://repos.influxdata.com/influxdb.key | sudo apt-key add -
      - echo "deb https://repos.influxdata.com/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/influxdb.list
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
  EOF
}
*/

/*

resource "digitalocean_domain" "main" {
  name = "bee-ci.pacia.tech"
}

resource "digitalocean_record" "frontend" {
  domain = digitalocean_domain.main.id
  type   = "CNAME"
  name   = "app"
  value  = "cname.vercel-dns.com."
  ttl    = 1800
}

resource "digitalocean_record" "backend" {
  domain = digitalocean_domain.main.id
  type   = "CNAME"
  name   = "backend"
  ttl    = 1800

  # It's a bit hard to get the value we need
  # value  = "bee-ci.pacia.tech."
  value = format("%s.", split("://", digitalocean_app.app.default_ingress)[1])
  # I also tried below, but it doesn't work. See also: https://github.com/digitalocean/terraform-provider-digitalocean/issues/1206
  # value  = format("%s.", digitalocean_app.app.live_domain)
  # value = format("%s.", digitalocean_app.app.live_domain)
}

*/

resource "digitalocean_domain" "default" {
  name = "bee-ci.karolak.cc"
}

resource "digitalocean_domain" "default2" {
  name = "beeci-backend.karolak.cc"
}
