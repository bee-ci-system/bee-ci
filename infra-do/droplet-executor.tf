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
      - echo "hello from cloud-init on executor! started! pwd is: $(pwd)" >> /root/hello.txt

    final_message: "executor ready!"
    EOF
}
