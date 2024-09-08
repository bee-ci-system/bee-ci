terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "2.40.0"
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
      type = "PRIMARY"
      zone = digitalocean_domain.main.name
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

      dynamic "env" {
        for_each = local.env_vars
        content {
          key   = env.value.key
          value = env.value.value
          scope = env.value.scope
          type  = lookup(env.value, "type", null)
        }
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
  name       = "bee"

  provisioner "local-exec" {
    command = <<EOF
      psql -f ../backend/sql-scripts/1-schema.sql
      psql -f ../backend/sql-scripts/2-triggers.sql
      psql -f ../backend/sql-scripts/3-views.sql
      psql -f ../backend/sql-scripts/100-seed.sql
    EOF

    environment = {
      "PGUSER"     = digitalocean_database_cluster.main-db-cluster.user # -U
      "PGHOST"     = digitalocean_database_cluster.main-db-cluster.host # -h
      "PGPORT"     = digitalocean_database_cluster.main-db-cluster.port # -p
      "PGDATABASE" = digitalocean_database_db.main-db.name              # -d
      "PGSSLMODE"  = "require"                                          # -c sslmode=require
      "PGPASSWORD" = digitalocean_database_cluster.main-db-cluster.password
    }
  }
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
  ttl    = 1800
}


resource "digitalocean_record" "backend" {
  domain = digitalocean_domain.main.id
  type   = "CNAME"
  name   = "backend"
  value  = format("%s.", split("://", digitalocean_app.app.default_ingress)[1])
  # value  = "bee-ci.pacia.tech."
  ttl = 1800
}

