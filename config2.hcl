storage "postgresql" {
  connection_url = "postgres://postgres:rails@localhost:5432/test?sslmode=disable"
  table="vault_kv_store"
}

listener "tcp" {
 address     = "127.0.0.1:8200"
 tls_disable = 1
}

"api_addr" = "http://127.0.0.1:8200"
"plugin_directory" ="/home/rails/GO/src/github.com/sethvargo/vault-secrets-gen"