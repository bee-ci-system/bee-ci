locals {
  image_id = "${var.region}-docker.pkg.dev/${random_string.project_id.result}/${google_artifact_registry_repository.default.repository_id}/bee-ci:latest"
}
