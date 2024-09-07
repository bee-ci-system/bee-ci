terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "2.40.0"
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
    digitalocean_database_cluster.main-db-cluster.urn,
    digitalocean_domain.main.urn,
  ]
}

resource "digitalocean_app" "app" {
  spec {
    name   = "bee-ci-tf"
    region = "sfo"

    domain {
      name = "backend.bee-ci.pacia.tech"
      type = "ALIAS"
      zone = "bee-ci.pacia.tech"
    }

    database {
      name         = digitalocean_database_db.main-db.name
      db_name      = digitalocean_database_db.main-db.name
      cluster_name = digitalocean_database_cluster.main-db-cluster.name
      production   = true
    }

    service {
      name               = "backend"
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "apps-s-1vcpu-0.5gb"

      env {
        key   = "PORT"
        value = "8080"
        scope = "RUN_TIME"
      }

      env {
        key   = "GITHUB_APP_ID"
        value = var.github_app_id
        scope = "RUN_TIME"
      }

      env {
        key   = "GITHUB_APP_CLIENT_ID"
        value = var.github_app_client_id
        scope = "RUN_TIME"
      }

      env {
        key   = "GITHUB_APP_WEBHOOK_SECRET"
        value = var.github_app_webhook_secret
        scope = "RUN_TIME"
        type  = "SECRET"
      }

      env {
        key   = "GITHUB_APP_PRIVATE_KEY_BASE64"
        value = var.github_app_private_key_base64
        scope = "RUN_TIME"
        type  = "SECRET"
      }

      env {
        key   = "GITHUB_APP_CLIENT_SECRET"
        value = var.github_app_client_secret
        scope = "RUN_TIME"
        type  = "SECRET"
      }

      env {
        key   = "DB_HOST"
        value = digitalocean_database_cluster.main-db-cluster.host
        scope = "RUN_TIME"
      }

      env {
        key   = "DB_PORT"
        value = digitalocean_database_cluster.main-db-cluster.port
        scope = "RUN_TIME"
      }

      env {
        key   = "DB_USER"
        value = digitalocean_database_cluster.main-db-cluster.user
        scope = "RUN_TIME"
      }

      env {
        key   = "DB_PASSWORD"
        value = digitalocean_database_cluster.main-db-cluster.password
        scope = "RUN_TIME"
        type  = "SECRET"
      }

      env {
        key   = "DB_NAME"
        value = digitalocean_database_db.main-db.name
        scope = "RUN_TIME"
      }

      env {
        key   = "DB_OPTS"
        value = "sslmode=require"
        scope = "RUN_TIME"
      }

      #dockerfile_path = "./backend/Dockerfile"
      # git {
      #   repo_clone_url = "https://github.com/bee-ci-system/bee-ci"
      #   branch         = "master"
      # }

      image {
        registry_type = "DOCR" # DigitalOcean Container Registry
        repository    = "backend"
        tag           = "latest"

        deploy_on_push {
          enabled = true
        }
      }

      health_check {
        http_path             = "/"
        initial_delay_seconds = 10
        period_seconds        = 5
        timeout_seconds       = 1
        success_threshold     = 3
        failure_threshold     = 3
      }
    }
  }
}

resource "digitalocean_database_db" "main-db" {
  cluster_id = digitalocean_database_cluster.main-db-cluster.id
  name       = "main-db"
}

resource "digitalocean_database_cluster" "main-db-cluster" {
  name       = "bee-postgres-cluster"
  engine     = "pg"
  version    = "16"
  size       = "db-s-1vcpu-1gb"
  region     = "sfo2"
  node_count = 1
}

resource "digitalocean_domain" "main" {
  name = "bee-ci.pacia.tech"
}

resource "digitalocean_record" "frontend" {
  domain = digitalocean_domain.main.id
  type   = "CNAME"
  name   = "app"
  value  = "cname.vercel-dns.com."
}


resource "digitalocean_record" "backend" {
  domain = digitalocean_domain.main.id
  type   = "CNAME"
  name   = "backend"
  value  = format("%s.", split("://", digitalocean_app.app.default_ingress)[1])
  # value  = "bee-ci.pacia.tech."
}

