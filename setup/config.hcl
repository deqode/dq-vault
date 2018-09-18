storage "postgresql" {
  connection_url = "postgres://vault:vault@database:5432/vault?sslmode=disable"
  table="vault_kv_store"
}

listener "tcp" {
 address     = "0.0.0.0:8200"
 tls_disable = 1
}

"api_addr" = "http://127.0.0.1:8200"
"plugin_directory" ="/vault/plugins"
"disable_mlock"=true