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
    digitalocean_droplet.executor.urn,
    digitalocean_droplet.influxdb.urn,
    digitalocean_volume.influxdb_volume.urn,
  ]
}

// General:
//
// CNAME bee-ci.pacia.tech      bee-ci-tf-rqd3k.ondigitalocean.app
//
// For Vercel:
//
// CNAME app.bee-ci.pacia.tech  cname.vercel-dns.com
// TXT   _vercel.pacia.tech        vc-domain-verify=app.bee-ci.pacia.tech,9f0455b9ebb2cffa2f86
//
// Ideally, we can get rid of Vercel, but it's blocked by:
// https://stackoverflow.com/questions/79029275/digital-ocean-app-platform-nextresponse-redirect-doesnt-work-status-code-2

resource "digitalocean_app" "app" {

  spec {
    name   = "bee-ci-tf"
    region = "sfo"

    // Requires the following record to be set on domain "pacia.tech"
    // Type  | Hostname          | Value
    // CNAME | bee-ci.pacia.tech | bee-ci-tf-rqd3k.ondigitalocean.app (or whatever DO generated)
    domain {
      name = "bee-ci.pacia.tech"
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
