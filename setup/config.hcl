storage "postgresql" {
  connection_url = "postgres://postgres:postgres@34.126.118.46:5432/vault?sslmode=disable"
  table="vault_kv_store"
}

listener "tcp" {
 address     = "34.126.118.46:8200"
 tls_disable = 1
}

"api_addr" = "http://34.126.118.46:8200"
"plugin_directory" ="/vault/plugins"
"disable_mlock"=true
