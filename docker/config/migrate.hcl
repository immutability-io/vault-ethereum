storage_source "postgresql" {
  connection_url = "postgres://tayto:bun_bo@10.10.0.106:2345/poc?sslmode=disable"
  table = "vault_kv_store"
  ha_table = "vault_ha_locks"
}

storage_destination "file" {
  path = "/Users/user/Downloads/kai/vault-ethereum/docker/config/datas"
}