disable_mlock = true
ui            = true

listener "tcp" {
  address     = "127.0.0.1:8200"
  tls_disable = "true"
}

storage "file" {
  path = "/home/deq/etc/tmp/vault-data"
}

plugin_directory = "/home/deq/etc/vault/plugins/vault_plugin"
api_addr = "http://127.0.0.1:8200"