"api_addr" = "http://127.0.0.1:8200"

storage "postgresql" {
  connection_url = "postgres://postgres:rails@localhost:5432/postgres?sslmode=disable"
  table = "mytable"	
}


listener "tcp" {
 address     = "127.0.0.1:8200"
 tls_disable = 1
}

"plugin_directory" = "/home/rails/blockchain/vault/ethereum-plugin"




