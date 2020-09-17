storage "postgresql" {
  connection_url = "postgres://vault:vault@database:5432/vault?sslmode=disable"
  table="vault_kv_store"
}

listener "tcp" {
 address     = "34.87.77.243"
 tls_disable = 1
}

"api_addr" = "http://34.87.77.243"
"plugin_directory" ="/vault/plugins"
"disable_mlock"=true
